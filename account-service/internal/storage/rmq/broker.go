package rmq

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/account-service/internal/config"
	"github.com/tehrelt/mu/account-service/internal/dto"
)

type Broker struct {
	cfg *config.Config
	ch  *amqp091.Channel
}

func New(cfg *config.Config, ch *amqp091.Channel) *Broker {
	return &Broker{cfg, ch}
}

func (b *Broker) PublishBalanceChanged(ctx context.Context, event *dto.EventBalanceChanged) error {
	j, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed marshal event data", slog.Any("event", event))
		return err
	}

	return b.ch.PublishWithContext(ctx, b.cfg.BalanceChanged.Exchange, b.cfg.BalanceChanged.Routing, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        j,
	})
}
