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
	BalanceChangedQueue      = "notification_service.balance.changed"
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

	balanceChangedMessages, err := c.manager.Consume(ctx, BalanceChangedQueue)
	if err != nil {
		return err
	}

	log := slog.With(slog.String("amqp", "amqp"))

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-ticketStatusChangedMessages:
			if err := c.handleTicketStatusChangedMessage(msg.Context(), msg); err != nil {
				log.Error("failed to handle message", slog.String("queue", TicketStatusChangedQueue), sl.Err(err))
			}
		case msg := <-balanceChangedMessages:
			if err := c.handleBalanceChangedEvent(msg.Context(), msg); err != nil {
				log.Error("failed to handle message", slog.String("queue", BalanceChangedQueue), sl.Err(err))
			}
		}
	}
}

func (c *AmqpConsumer) handleBalanceChangedEvent(ctx context.Context, msg *rmqmanager.TracedDelivery) (err error) {

	fn := "handleBalanceChangedEvent"
	log := slog.With(sl.Method(fn))

	defer func() {
		if err != nil {
			return
		}

		err = msg.Ack(false)
		if err != nil {
			log.Error("failed to ack message", slog.String("queue", BalanceChangedQueue), sl.Err(err))
		}
	}()

	event := &events.IncomingBalanceChanged{}
	if err := json.Unmarshal(msg.Body, event); err != nil {
		log.Error("failed to unmarshal message", sl.Err(err))
		return err
	}

	if err := c.uc.HandleBalanceChangedEvent(ctx, event); err != nil {
		log.Error("failed to handle balance changed event", sl.Err(err))
		return err
	}

	return nil
}
