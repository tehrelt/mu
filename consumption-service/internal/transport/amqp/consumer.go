package amqp

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/internal/config"
	"github.com/tehrelt/mu/consumption-service/internal/usecase"
)

type AmqpConsumer struct {
	cfg     *config.Config
	manager *rmqmanager.RabbitMqManager
	uc      *usecase.UseCase
}

func New(
	cfg *config.Config,
	ch *amqp091.Channel,
	uc *usecase.UseCase,
) *AmqpConsumer {
	return &AmqpConsumer{
		cfg:     cfg,
		manager: rmqmanager.New(ch),
		uc:      uc,
	}
}

func (c *AmqpConsumer) Run(ctx context.Context) error {

	serviceConnectedMessages, err := c.manager.Consume(ctx, c.cfg.ConnectServiceExchange.Queue)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-serviceConnectedMessages:
			if err := c.handleServiceConnectedEvent(msg.Context(), msg); err != nil {
				slog.Error("failed to handle service connected event", sl.Err(err))
			}
		}

	}

}
