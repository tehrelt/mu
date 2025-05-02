package models

import (
	"time"

	"github.com/google/uuid"
)

type Cabinet struct {
	Id        uuid.UUID
	AccountId uuid.UUID
	ServiceId uuid.UUID
	Consumed  uint64
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type ConsumptionLog struct {
	Id            uuid.UUID
	Consumed      uint64
	CabinetId     uuid.UUID
	AccountId     uuid.UUID
	ServiceId     uuid.UUID
	PaymentId     uuid.UUID
	ConsumptionId uuid.UUID
	CreatedAt     time.Time
}

type Service struct {
	Id   uuid.UUID
	Name string
	Rate uint64
	Unit string
}

type Account struct {
	Id uuid.UUID
}

type Charge struct {
	PaymentId uuid.UUID
	Amount    uint64
	Message   string
}
