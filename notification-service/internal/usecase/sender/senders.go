package sender

import (
	"context"

	"github.com/tehrelt/mu/notification-service/internal/events"
	"github.com/tehrelt/mu/notification-service/internal/storage/rmq"
)

type telegram struct {
}

var _ Sender = (*telegram)(nil)

func (s *telegram) send(ctx context.Context, broker *rmq.Broker, event events.Event) error {
	return broker.SendTelegramNotification(ctx, event)
}

type email struct {
}

var _ Sender = (*email)(nil)

func (s *email) send(ctx context.Context, broker *rmq.Broker, event events.Event) error {
	return broker.SendEmailNotification(ctx, event)
}
