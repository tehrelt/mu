package accountstorage

import "github.com/jmoiron/sqlx"

type AccountStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *AccountStorage {
	return &AccountStorage{db: db}
}
