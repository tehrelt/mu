package storage

import "errors"

var (
	ErrInvalidId         error = errors.New("invalid id")
	ErrTicketNotFound          = errors.New("ticket not found")
	ErrInvalidTicketType       = errors.New("invalid ticket type")
	ErrNoIdProvided            = errors.New("no id provided")
)
