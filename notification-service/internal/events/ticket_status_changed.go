package events

type IncomingTicketStatusChanged struct {
	TicketId  string `json:"ticketId"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type TicketStatusChanged struct {
	EventHeader
	TicketId  string `json:"ticketId"`
	NewStatus string `json:"newStatus"`
}

func NewTicketStatusChanged(header EventHeader, ticketId string, newStatus string) TicketStatusChanged {
	return TicketStatusChanged{
		EventHeader: header,
		TicketId:    ticketId,
		NewStatus:   newStatus,
	}
}
