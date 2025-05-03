package amqp

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu/housing-service/internal/config"
	"github.com/tehrelt/mu/housing-service/internal/storage/pg/housestorage"
	ratepb "github.com/tehrelt/mu/housing-service/pkg/pb/ratespb"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
)

type AmqpConsumer struct {
	cfg     *config.Config
	manager *rmqmanager.RabbitMqManager
	storage *housestorage.HouseStorage
	rateapi ratepb.RateServiceClient
}

func New(cfg *config.Config, ch *amqp091.Channel, s *housestorage.HouseStorage, rapi ratepb.RateServiceClient) *AmqpConsumer {
	return &AmqpConsumer{
		cfg:     cfg,
		storage: s,
		rateapi: rapi,
		manager: rmqmanager.New(ch),
	}
}

func (c *AmqpConsumer) Run(ctx context.Context) error {

	connectServiceQueue, err := c.manager.Consume(ctx, c.cfg.ConnectServiceExchange.Queue)
	if err != nil {
		slog.Error("failed to consume connect service event", sl.Err(err))
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-connectServiceQueue:
			if err := c.handleConnectServiceEvent(msg.Context(), msg); err != nil {
				slog.Error("failed to consume connect service event", sl.Err(err))
			}
		}
	}
}
