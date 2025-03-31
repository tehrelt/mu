package amqp

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/account-service/internal/config"
	"github.com/tehrelt/mu/account-service/internal/storage/pg/accountstorage"
	"github.com/tehrelt/mu/account-service/internal/storage/rmq"
	"github.com/tehrelt/mu/account-service/pkg/pb/billingpb"
)

type AmqpConsumer struct {
	cfg        *config.Config
	manager    *rmqmanager.RabbitMqManager
	storage    *accountstorage.AccountStorage
	broker     *rmq.Broker
	billingApi billingpb.BillingServiceClient
}

func New(
	cfg *config.Config,
	ch *amqp091.Channel,
	s *accountstorage.AccountStorage,
	b *rmq.Broker,
	billingApi billingpb.BillingServiceClient,
) *AmqpConsumer {
	return &AmqpConsumer{
		cfg:        cfg,
		manager:    rmqmanager.New(ch),
		storage:    s,
		broker:     b,
		billingApi: billingApi,
	}
}

func (c *AmqpConsumer) Run(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := c.handleEvents(ctx); err != nil {
					slog.Error("failed to consume connec service event", sl.Err(err))
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func (c *AmqpConsumer) handleEvents(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := c.manager.Consume(ctx, c.cfg.PaymentStatusChanged.Routing, c.handlePaymentStatusChangedEvent); err != nil {
				slog.Error("failed to consume payment status changed event", sl.Err(err))
			}
		}
	}
}
