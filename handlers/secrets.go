package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"www.github.com/mrgne1/paperhat/encryption"

	"github.com/google/uuid"
)

var ErrNoSecretFound error = errors.New("No secret found")

type CreateSecretResponse struct {
	Id         string `json:"id"`
	Key        string `json:"key"`
}

type Secret struct {
	Id        uuid.UUID
	Value     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

var maxDuration int64 = 86400
var minDuration int64 = 60

func (c *ApiConfig) CreateSecretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("CreateSecretHandler Error: %v\n", err)
			secretHandlerError(w, "Error reading secret from HTTP message")
			return
		}

		id, err := uuid.NewRandom()
		if err != nil {
			log.Printf("CreateSecretHandler Error: %v\n", err)
			secretHandlerError(w, "Error creating the secret")
			return
		}

		cipherText, keyText, err := encryption.Encrypt(string(body))
		if err != nil {
			log.Printf("CreateSecretHandler Error: %v\n", err)
			secretHandlerError(w, "Error Encrypting secret")
			return
		}

		durationString := r.URL.Query().Get("duration")
		duration, err := strconv.ParseInt(durationString, 0, 32)
		if err != nil {
			log.Printf("CreateSecretHandler Invalid Duration %v\n", durationString)
			duration = minDuration
		} else if duration <= minDuration || duration > maxDuration  {
			log.Printf("CreateSecretHandler Duration out of bounds %v\n", duration)
			duration = minDuration
		}

		secret := Secret{
			Id:        id,
			Value:     cipherText,
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(time.Duration(duration) * time.Second),
		}

		err = c.createSecret(secret)
		if err != nil {
			log.Printf("CreateSecretHandler Error: %v\n", err)
			secretHandlerError(w, "Error creating the secret in DB")
			return
		}

		log.Printf("New secret created! Expires at %v\n", secret.ExpiresAt)

		resp := CreateSecretResponse{
			Id:         secret.Id.String(),
			Key:        keyText,
		}

		respBody, _ := json.Marshal(resp)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(respBody)
	})
}

func (c *ApiConfig) ReadSecretHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keyText := r.PathValue("keyText")
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			log.Printf("ReadSecretHandler Error: %v\n", err)
			secretHandlerError(w, "Error parsing the secret id")
			return
		}

		secret, err := c.readSecret(id)
		if err != nil {
			log.Printf("ReadSecretHandler Error: %v\n", err)
			secretHandlerError(w, "Error Reading secret from DB")
			return
		}

		err = c.deleteSecret(id)
		if err != nil {
			log.Printf("ReadSecretHandler Error: %v\n", err)
			secretHandlerError(w, "Error Reading secret from DB")
			return
		}

		value, err := encryption.Decrypt(secret.Value, keyText)
		if err != nil {
			log.Printf("ReadSecretHandler Error: %v\n", err)
			secretHandlerError(w, "Error decrypting secret")
			return
		}

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(value))
	})
}

func secretHandlerError(w http.ResponseWriter, message string) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(500)
	w.Write([]byte(message))
}

var createSecretSql string = `
insert into secrets (id, value, created_at, expires_at)
values (?,?,?,?);
`

func (c *ApiConfig) createSecret(secret Secret) error {
	_, err := c.secrets.ExecContext(
		context.Background(),
		createSecretSql,
		secret.Id,
		secret.Value,
		secret.CreatedAt,
		secret.ExpiresAt,
	)
	return err
}

var readSecretSql string = `
select id, value, created_at, expires_at
from secrets
where secrets.id = ?;
`

func (c *ApiConfig) readSecret(id uuid.UUID) (Secret, error) {
	row := c.secrets.QueryRowContext(
		context.Background(),
		readSecretSql,
		id,
	)
	var secret Secret
	err := row.Scan(
		&secret.Id,
		&secret.Value,
		&secret.CreatedAt,
		&secret.ExpiresAt,
	)
	return secret, err
}

var deleteSecretSql string = `
delete from secrets
where id = ?;
`

func (c *ApiConfig) deleteSecret(id uuid.UUID) error {
	_, err := c.secrets.ExecContext(
		context.Background(),
		deleteSecretSql,
		id,
	)
	return err
}

var deleteExpiredSql string = `
delete from secrets
where expires_at < ?;
`

func (cfg *ApiConfig) DeleteExpired() error {
	result, err := cfg.secrets.ExecContext(
		context.Background(),
		deleteExpiredSql,
		time.Now().UTC(),
	)

	n, e := result.RowsAffected()
	if e != nil {
		log.Printf("deleteExpiredRows Error: %v\n", e)
	} else {
		log.Printf("Deleted %v expired rows\n", n)
	}

	return err
}
