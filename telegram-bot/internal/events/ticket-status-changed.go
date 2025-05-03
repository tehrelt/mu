package events

type TicketStatusChanged struct {
	EventHeader
	TicketId  string `json:"ticketId"`
	NewStatus string `json:"newStatus"`
}
