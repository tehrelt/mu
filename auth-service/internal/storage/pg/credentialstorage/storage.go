package credentialstorage

import "github.com/jmoiron/sqlx"

const traceKey = "credentialstorage"

type CredentialStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *CredentialStorage {
	return &CredentialStorage{db}
}
