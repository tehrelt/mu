package sender

import (
	"context"

	"github.com/tehrelt/mu/notification-service/internal/dto"
	"github.com/tehrelt/mu/notification-service/internal/events"
	"github.com/tehrelt/mu/notification-service/internal/storage/rmq"
)

type Sender interface {
	send(ctx context.Context, broker *rmq.Broker, event events.Event) error
}

type manager struct {
	senders []Sender
	broker  *rmq.Broker
}

func (m *manager) appendSender(s Sender) {
	m.senders = append(m.senders, s)
}

func (m *manager) Send(ctx context.Context, event events.Event) error {
	for _, sender := range m.senders {
		if err := sender.send(ctx, m.broker, event); err != nil {
			return err
		}
	}
	return nil
}

func New(broker *rmq.Broker, settings *dto.UserSettings) *manager {
	m := &manager{
		broker: broker,
	}

	m.build(settings)

	return m
}

func (m *manager) build(settings *dto.UserSettings) {
	opts := make([]SenderOption, 0)

	opts = append(opts, WithSender(&email{}))

	if settings.TelegramChatId != nil {
		opts = append(opts, WithSender(&telegram{}))
	}

	for _, opt := range opts {
		opt(m)
	}
}
