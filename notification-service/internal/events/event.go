package events

import (
	"time"

	"github.com/tehrelt/mu/notification-service/internal/dto"
)

type EventHeader struct {
	EventType EventType         `json:"type"`
	UserId    string            `json:"userId"`
	Timestamp time.Time         `json:"timestamp"`
	Settings  *dto.UserSettings `json:"settings"`
}

func NewEventHeader(eventType EventType, userId string) EventHeader {
	return EventHeader{
		EventType: eventType,
		UserId:    userId,
		Timestamp: time.Now(),
	}
}

func (h *EventHeader) Header() *EventHeader {
	return h
}

func (h *EventHeader) Type() EventType {
	return h.EventType
}

func (h *EventHeader) SetSettings(s *dto.UserSettings) {
	h.Settings = s
}

type Event interface {
	Header() *EventHeader
	Type() EventType
	SetSettings(*dto.UserSettings)
}
