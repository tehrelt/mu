package rmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu/ticket-service/internal/config"
	"github.com/tehrelt/mu/ticket-service/internal/dto"
)

const traceKey = "rmq broker"

type Broker struct {
	cfg     *config.Config
	manager *rmqmanager.RabbitMqManager
}

func New(cfg *config.Config, ch *amqp091.Channel) *Broker {
	return &Broker{
		cfg:     cfg,
		manager: rmqmanager.New(ch),
	}
}

func (b *Broker) PublishStatusNewAccount(ctx context.Context, event *dto.EventTicketStatusChanged) error {
	event.Timestamp = time.Now()

	j, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed marshal event data", slog.Any("event", event))
		return err
	}

	exchange := b.cfg.TicketStatusChangedQueue.Exchange
	rk := b.cfg.TicketStatusChangedQueue.NewAccountRoute

	return b.publish(ctx, exchange, rk, j)
}

func (b *Broker) PublishStatusConnectService(ctx context.Context, event *dto.EventTicketStatusChanged) error {
	event.Timestamp = time.Now()

	j, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed marshal event data", slog.Any("event", event))
		return err
	}

	exchange := b.cfg.TicketStatusChangedQueue.Exchange
	rk := b.cfg.TicketStatusChangedQueue.ConnectServiceRoute

	return b.publish(ctx, exchange, rk, j)
}

func (b *Broker) publish(ctx context.Context, exchange, key string, data []byte) error {
	slog.Info(
		"sending event",
		slog.String("exchange", exchange),
		slog.String("routingKey", key),
		slog.Any("event", data),
	)

	return b.manager.Publish(ctx, exchange, key, data)
}
