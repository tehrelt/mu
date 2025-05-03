package events

type EventType string

const (
	EventTicketStatusChanged EventType = "ticket_status_changed"
	EventBalanceChanged      EventType = "balance_changed"
	EventRateChanged         EventType = "rate_changed"
)
