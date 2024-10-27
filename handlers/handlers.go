package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "modernc.org/sqlite"
)

var ErrOpeningSecretDB error = errors.New("Can't Open secrets database")

type ApiConfig struct {
	Version string
	secrets *sql.DB
	hostUrl string
}

func NewApiConfig(hostUrl string) (ApiConfig, error) {
	secrets, err := getSecretDb()
	if err != nil {
		log.Println(err)
		return ApiConfig{}, ErrOpeningSecretDB
	}

	return ApiConfig{
			Version: "1",
			secrets: secrets,
			hostUrl: hostUrl,
		},
		nil
}

func getSecretDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./secrets.db")
	if err != nil {
		return &sql.DB{}, err
	}

	_, err = db.ExecContext(
		context.Background(),
		`create table if not exists secrets (
			id uuid primary key,
			value string,
			created_at timestamp not null,
			expires_at timestamp not null
		)`,
	)
	if err != nil {
		return &sql.DB{}, err
	}

	return db, nil
}

