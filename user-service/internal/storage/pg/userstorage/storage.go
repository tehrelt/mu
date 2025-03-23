package userstorage

import "github.com/jmoiron/sqlx"

type UserStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}
