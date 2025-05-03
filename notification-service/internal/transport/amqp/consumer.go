package amqp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/config"
	"github.com/tehrelt/mu/notification-service/internal/events"
	"github.com/tehrelt/mu/notification-service/internal/storage/rmq"
	"github.com/tehrelt/mu/notification-service/internal/usecase"
)

const (
	TicketStatusChangedQueue = "notification_service.ticket.status_changed"
)

type AmqpConsumer struct {
	cfg     *config.Config
	manager *rmqmanager.RabbitMqManager
	broker  *rmq.Broker
	uc      *usecase.UseCase
}

func New(
	cfg *config.Config,
	ch *amqp091.Channel,
	b *rmq.Broker,
	uc *usecase.UseCase,
) *AmqpConsumer {
	return &AmqpConsumer{
		cfg:     cfg,
		manager: rmqmanager.New(ch),
		broker:  b,
		uc:      uc,
	}
}

func (c *AmqpConsumer) Run(ctx context.Context) error {

	ticketStatusChangedMessages, err := c.manager.Consume(ctx, TicketStatusChangedQueue)
	if err != nil {
		return err
	}

	log := slog.With(slog.String("amqp", "amqp"))

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-ticketStatusChangedMessages:
			if err := c.handleTicketStatusChangedMessage(ctx, msg); err != nil {
				log.Error("failed to handle message", slog.String("queue", TicketStatusChangedQueue), sl.Err(err))
			}
		}
	}
}

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
