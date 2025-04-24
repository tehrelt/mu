package rmq

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/housing-service/internal/config"
	"github.com/tehrelt/mu/housing-service/internal/dto"
)

type Broker struct {
	cfg *config.Config
	ch  *amqp091.Channel
}

func New(cfg *config.Config, ch *amqp091.Channel) *Broker {
	return &Broker{cfg, ch}
}

func (b *Broker) PublishServiceConnectedEvent(ctx context.Context, event *dto.EventServiceConnected) error {
	j, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed marshal event data", slog.Any("event", event))
		return err
	}

	return b.ch.PublishWithContext(ctx, b.cfg.ServiceConnectedQueue.Exchange, b.cfg.ServiceConnectedQueue.Routing, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        j,
	})
}
