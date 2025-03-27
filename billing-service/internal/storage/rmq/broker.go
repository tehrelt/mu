package rmq

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/billing-service/internal/config"
	"github.com/tehrelt/mu/billing-service/internal/dto"
)

type Broker struct {
	cfg *config.Config
	ch  *amqp091.Channel
}

func New(cfg *config.Config, ch *amqp091.Channel) *Broker {
	return &Broker{cfg, ch}
}

func (b *Broker) PublishStatusChanged(ctx context.Context, event *dto.EventStatusChanged) error {
	j, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed marshal event data", slog.Any("event", event))
		return err
	}

	return b.ch.PublishWithContext(ctx, b.cfg.PaymentStatusChanged.Exchange, b.cfg.PaymentStatusChanged.Routing, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        j,
	})
}
