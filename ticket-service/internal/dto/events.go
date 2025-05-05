package dto

import (
	"time"

	"github.com/tehrelt/mu/ticket-service/internal/models"
)

type EventTicketStatusChanged struct {
	TicketId  string              `json:"ticketId"`
	Ticket    models.Ticket       `json:"ticket"`
	Status    models.TicketStatus `json:"status"`
	Timestamp time.Time           `json:"timestamp"`
}
