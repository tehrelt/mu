package events

type TicketStatusChanged struct {
	EventHeader
	TicketId  string `json:"ticketId"`
	NewStatus string `json:"newStatus"`
}

type BalanceChanged struct {
	EventHeader
	OldBalance int64  `json:"oldBalance"`
	NewBalance int64  `json:"newBalance"`
	Address    string `json:"address"`
	Reason     string `json:"reason"`
}

func (e *BalanceChanged) Delta() int64 {
	return e.NewBalance - e.OldBalance
}
