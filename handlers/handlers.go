package handlers

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

var ErrOpeningSecretDB error = errors.New("Can't Open secrets database")

type ApiConfig struct {
	Version string
	secrets *sql.DB
}

func NewApiConfig(dbPath string) (ApiConfig, error) {
	secrets, err := getSecretDb(dbPath)
	if err != nil {
		log.Println(err)
		return ApiConfig{}, ErrOpeningSecretDB
	}

	return ApiConfig{
			Version: "1",
			secrets: secrets,
		},
		nil
}

func getSecretDb(dbPath string) (*sql.DB, error) {
	var dbType string
	if strings.HasPrefix(dbPath, "libsql://") {
		dbType = "libsql"
	} else {
		dbType = "sqlite"
	}

	db, err := sql.Open(dbType, dbPath)
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

