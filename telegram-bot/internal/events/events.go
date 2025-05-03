package events

import "time"

type EventType string

const (
	EventTicketStatusChanged EventType = "ticket_status_changed"
	EventBalanceChanged      EventType = "balance_changed"
	EventRateChanged         EventType = "rate_changed"
)

type EventHeader struct {
	EventType EventType    `json:"type"`
	UserId    string       `json:"userId"`
	Timestamp time.Time    `json:"timestamp"`
	Settings  UserSettings `json:"settings"`
}

func (eh *EventHeader) Header() *EventHeader {
	return eh
}

type UserSettings struct {
	Email          string `json:"email"`
	TelegramChatId string `json:"telegramChatId"`
}

type Event interface {
	Header() *EventHeader
}
