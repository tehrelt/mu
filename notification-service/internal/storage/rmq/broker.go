package rmq

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu/notification-service/internal/config"
	"github.com/tehrelt/mu/notification-service/internal/events"
)

const (
	telegramRK = "telegram"
	emailRK    = "email"
)

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

func (b *Broker) sendNotification(ctx context.Context, rk string, event events.Event) error {

	j, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return b.manager.Publish(ctx, b.cfg.NotificationSendExchange.Exchange, rk, j)
}

func (b *Broker) SendTelegramNotification(ctx context.Context, event events.Event) error {
	return b.sendNotification(ctx, telegramRK, event)
}

func (b *Broker) SendEmailNotification(ctx context.Context, event events.Event) error {
	return b.sendNotification(ctx, emailRK, event)
}
