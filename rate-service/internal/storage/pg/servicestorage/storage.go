package servicestorage

import (
	"github.com/jmoiron/sqlx"
)

type ServiceStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *ServiceStorage {
	return &ServiceStorage{
		db: db,
	}
}
