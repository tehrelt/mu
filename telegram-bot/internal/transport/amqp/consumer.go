package amqp

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/telegram-bot/internal/config"
	"github.com/tehrelt/mu/telegram-bot/internal/usecase"
)

const (
	NotifcationQueue = "telegram_bot.notifications"
)

type Consumer struct {
	cfg     *config.Config
	manager *rmqmanager.RabbitMqManager
	uc      *usecase.UseCase
	logger  *slog.Logger
}

func New(cfg *config.Config, uc *usecase.UseCase) (*Consumer, error) {
	c := &Consumer{
		cfg:    cfg,
		uc:     uc,
		logger: slog.With(sl.Module("RMQ Consumer")),
	}

	if err := c.setup(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Consumer) setup() error {

	conn, err := amqp091.Dial(c.cfg.AMQP.ConnectionString())
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(NotifcationQueue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := ch.QueueBind(queue.Name, "telegram", c.cfg.NotificationSendExchange, false, nil); err != nil {
		return err
	}

	c.manager = rmqmanager.New(ch)

	return nil
}

func (c *Consumer) Run(ctx context.Context) error {
	messages, err := c.manager.Consume(ctx, NotifcationQueue)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-messages:
			if err := c.handleMessage(ctx, msg); err != nil {
				c.logger.Error("failed to handle message", sl.Err(err))
			}
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg *rmqmanager.TracedDelivery) (err error) {

	defer func() {
		if err != nil {
			return
		}

	}()

	body := msg.Body

	c.logger.Info("event ", slog.String("body", string(body)))

	return nil
}
