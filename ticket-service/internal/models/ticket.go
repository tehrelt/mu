package models

import (
	"errors"
	"time"
)

type Ticket interface {
	Header() *TicketHeader
}

var (
	ErrInvalidTicketType = errors.New("invalid ticket type")
)

type TicketHeader struct {
	Id         string       `bson:"_id"`
	TicketType TicketType   `bson:"type"`
	Status     TicketStatus `bson:"status"`
	CreatedBy  string       `bson:"user_id"`
	CreatedAt  time.Time    `bson:"createdAt"`
	UpdatedAt  time.Time    `bson:"updatedAt"`
}

func (t *TicketHeader) Header() *TicketHeader {
	return t
}

func (t *TicketHeader) SetId(id string) {
	t.Id = id
}

type NewAccountTicket struct {
	TicketHeader
	UserId  string `bson:"userId"`
	Address string `bson:"address"`
}

type ConnectServiceTicket struct {
	TicketHeader
	ServiceId string `bson:"serviceId"`
	AccountId string `bson:"accountId"`
}
