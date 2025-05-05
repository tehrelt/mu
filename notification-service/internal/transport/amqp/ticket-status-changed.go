package amqp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu/notification-service/internal/events"
)

func (c *AmqpConsumer) handleTicketStatusChangedMessage(ctx context.Context, msg *rmqmanager.TracedDelivery) (err error) {
	defer func() {
		if err != nil {
			msg.Reject(false)
		} else {
			err = msg.Ack(false)
		}
	}()

	log := slog.With(slog.String("queue_handler", TicketStatusChangedQueue))
	log.Info("incoming message", slog.String("body", string(msg.Body)))

	event := &events.IncomingTicketStatusChanged{}
	if err := json.Unmarshal(msg.Body, event); err != nil {
		return err
	}

	if err := c.uc.HandleTicketStatusChangedEvent(ctx, event); err != nil {
		return err
	}

	return nil
}
