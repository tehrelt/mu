package housestorage

import "github.com/jmoiron/sqlx"

type HouseStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *HouseStorage {
	return &HouseStorage{db: db}
}
