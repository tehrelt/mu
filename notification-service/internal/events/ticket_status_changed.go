package events

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
