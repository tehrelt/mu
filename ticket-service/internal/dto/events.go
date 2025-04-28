package dto

import "time"

type EventTicketStatusChanged struct {
	TicketId  string    `json:"ticketId"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
