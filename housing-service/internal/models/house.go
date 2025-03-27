package models

import (
	"time"

	"github.com/google/uuid"
)

type House struct {
	Id                uuid.UUID
	RoomsQuantity     int
	ResidentsQuantity int
	Address           string
	ConnectedServices []uuid.UUID
	CreatedAt         time.Time
	UpdatedAt         *time.Time
}
