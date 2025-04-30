package dto

import (
	"github.com/google/uuid"
)

type NewCabinet struct {
	AccountId uuid.UUID
	ServiceId uuid.UUID
}

type NewConsume struct {
	FindCabinet
	Consumed uint64
}

type ConsumeCreated struct {
	LogId     uuid.UUID
	PaymentId uuid.UUID
	Amount    uint64
}

type NewConsumeLog struct {
	Consumed  uint64
	PaymentId uuid.UUID
	CabinetId uuid.UUID
}

type FindCabinet struct {
	Id        uuid.UUID
	AccountId uuid.UUID
	ServiceId uuid.UUID
}

type Charge struct {
	AccountId uuid.UUID
	ServiceId uuid.UUID
	Amount    uint64
}

type UpdateCabinet struct {
	Id            uuid.UUID
	ConsumedDelta uint64
}
