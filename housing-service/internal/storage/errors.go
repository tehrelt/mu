package storage

import "errors"

var (
	ErrHouseNotFound        error = errors.New("house not found")
	ErrHouseAlreadyAcquired error = errors.New("house already acquired")
)
