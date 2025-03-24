package rolestorage

import "github.com/jmoiron/sqlx"

type RoleStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *RoleStorage {
	return &RoleStorage{db}
}
