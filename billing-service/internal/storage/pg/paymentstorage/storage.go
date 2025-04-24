package paymentstorage

import (
	"github.com/jmoiron/sqlx"
)

const (
	traceKey = "payment-storage"
)

type PaymentStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PaymentStorage {
	return &PaymentStorage{db: db}
}
