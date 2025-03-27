package paymentstorage

import "github.com/jmoiron/sqlx"

type PaymentStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PaymentStorage {
	return &PaymentStorage{db: db}
}
