package amqp

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/housing-service/internal/config"
	"github.com/tehrelt/mu/housing-service/internal/storage/pg/housestorage"
	"github.com/tehrelt/mu/housing-service/internal/storage/rmq"
	ratepb "github.com/tehrelt/mu/housing-service/pkg/pb/ratespb"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
)

type AmqpConsumer struct {
	cfg     *config.Config
	channel *amqp091.Channel
	storage *housestorage.HouseStorage
	broker  *rmq.Broker
	rateapi ratepb.RateServiceClient
}

func New(cfg *config.Config, ch *amqp091.Channel, s *housestorage.HouseStorage, b *rmq.Broker, rapi ratepb.RateServiceClient) *AmqpConsumer {
	return &AmqpConsumer{
		cfg:     cfg,
		channel: ch,
		storage: s,
		broker:  b,
		rateapi: rapi,
	}
}

func (c *AmqpConsumer) Run(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := c.ConsumeConnectServiceEvent(ctx); err != nil {
					slog.Error("failed to consume connec service event", sl.Err(err))
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}
