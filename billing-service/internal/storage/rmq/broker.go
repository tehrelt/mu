package rmq

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu/billing-service/internal/config"
	"github.com/tehrelt/mu/billing-service/internal/dto"
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

func (b *Broker) PublishStatusChanged(ctx context.Context, event *dto.EventStatusChanged) error {

	fn := "broker.PublishStatusChanged"
	log := slog.With(slog.String("fn", fn))

	event.Timestamp = time.Now()

	j, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed marshal event data", slog.Any("event", event))
		return err
	}

	log.Info(
		"sending event",
		slog.String("exchange", b.cfg.PaymentStatusChanged.Exchange),
		slog.String("routing", b.cfg.PaymentStatusChanged.Routing),
		slog.Any("event", event),
	)

	return b.manager.Publish(ctx, b.cfg.PaymentStatusChanged.Exchange, b.cfg.PaymentStatusChanged.Routing, j)
}
