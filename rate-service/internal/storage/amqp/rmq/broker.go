package rmq

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/rate-service/internal/config"
	"github.com/tehrelt/mu/rate-service/internal/models"
)

type RabbitMq struct {
	cfg     *config.Config
	channel *amqp.Channel
	q       *amqp.Queue
}

func New(cfg *config.Config, channel *amqp.Channel, q *amqp.Queue) *RabbitMq {
	return &RabbitMq{
		cfg:     cfg,
		channel: channel,
		q:       q,
	}
}

func (r *RabbitMq) NotifyRateChanged(ctx context.Context, event *models.EventRateChanged) error {
	j, err := json.Marshal(event)
	if err != nil {
		return err
	}

	slog.Info("publishing event rate changed", slog.Any("event", event))

	return r.channel.PublishWithContext(ctx, r.cfg.Amqp.RateChangedExchange, r.cfg.Amqp.RateChangedRoutingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        j,
	})
}
