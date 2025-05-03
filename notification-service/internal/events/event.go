package events

import "time"

type EventHeader struct {
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

func NewEventHeader(eventType EventType) EventHeader {
	return EventHeader{
		Type:      eventType,
		Timestamp: time.Now(),
	}
}

type Event interface {
	Type() EventType
}
