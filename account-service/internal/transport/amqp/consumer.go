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
	"github.com/tehrelt/mu/account-service/pkg/pb/housepb"
	"github.com/tehrelt/mu/account-service/pkg/pb/ticketpb"
)

type AmqpConsumer struct {
	cfg        *config.Config
	manager    *rmqmanager.RabbitMqManager
	storage    *accountstorage.AccountStorage
	broker     *rmq.Broker
	houseApi   housepb.HouseServiceClient
	billingApi billingpb.BillingServiceClient
	ticketApi  ticketpb.TicketServiceClient
}

func New(
	cfg *config.Config,
	ch *amqp091.Channel,
	s *accountstorage.AccountStorage,
	b *rmq.Broker,
	houseApi housepb.HouseServiceClient,
	billingApi billingpb.BillingServiceClient,
	ticketApi ticketpb.TicketServiceClient,
) *AmqpConsumer {
	return &AmqpConsumer{
		cfg:        cfg,
		manager:    rmqmanager.New(ch),
		storage:    s,
		broker:     b,
		houseApi:   houseApi,
		billingApi: billingApi,
		ticketApi:  ticketApi,
	}
}

func (c *AmqpConsumer) Run(ctx context.Context) error {

	paymentStatusChangedQueue, err := c.manager.Consume(ctx, c.cfg.PaymentStatusChanged.Routing)
	if err != nil {
		slog.Error("failed to consume payment status changed event", sl.Err(err))
		return err
	}

	newAccountQueue, err := c.manager.Consume(ctx, c.cfg.TicketStatusChanged.NewAccountRoute)
	if err != nil {
		slog.Error("failed to consume ticket status changed event", sl.Err(err))
		return err
	}

	ConnectServiceQueue, err := c.manager.Consume(ctx, c.cfg.TicketStatusChanged.ConnectServiceRoute)
	if err != nil {
		slog.Error("failed to consume ticket status changed event", sl.Err(err))
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-paymentStatusChangedQueue:
			if err := c.handlePaymentStatusChangedEvent(msg.Context(), msg); err != nil {
				slog.Error("failed to consume payment status changed event", sl.Err(err))
			}
		case msg := <-newAccountQueue:
			if err := c.handleNewAccountEvent(msg.Context(), msg); err != nil {
				slog.Error("failed to consume ticket status changed event", sl.Err(err))
			}
		case msg := <-ConnectServiceQueue:
			if err := c.handleConnectServiceEvent(msg.Context(), msg); err != nil {
				slog.Error("failed to consume ticket status changed event", sl.Err(err))
			}
		}
	}
}
