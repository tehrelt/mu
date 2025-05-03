package storage

import "errors"

var (
	ErrUserNotFound      error = errors.New("user not found")
	ErrUserAlreadyExists error = errors.New("user already exists")
)
