package services

import "errors"

var (
	ErrInvalidCredentials error = errors.New("invalid credentials")
	ErrInvalidSession     error = errors.New("invalid session")
	ErrInvalidToken       error = errors.New("invalid token")
	ErrUserExists         error = errors.New("user exists")
	ErrForbidden          error = errors.New("forbidden")
)
