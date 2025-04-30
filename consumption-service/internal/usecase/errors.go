package usecase

import "errors"

var (
	ErrConsumptionNotFound = errors.New("consumption not found")
	ErrServiceNotFound     = errors.New("service not found")
	ErrAccountNotFound     = errors.New("account not found")

	ErrPaymentServiceUnavailable = errors.New("payment service unavailable")
)
