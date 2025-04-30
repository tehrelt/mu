package amqp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
)

func (c *AmqpConsumer) ConsumeServiceConnected(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := c.manager.Consume(ctx, c.cfg.ServiceConnectedQueue.Routing, c.handleServiceConnectedEvent); err != nil {
				slog.Error("failed to consume", sl.Err(err))
				return err
			}
		}
	}
}

func (c *AmqpConsumer) handleServiceConnectedEvent(ctx context.Context, msg amqp091.Delivery) (err error) {

	defer func() {
		if err != nil {
			msg.Nack(false, true)
			return
		}

		err = msg.Ack(false)
	}()

	body := msg.Body
	event := dto.EventServiceConnected{}
	if err := json.Unmarshal(body, &event); err != nil {
		slog.Error("failed to unmarshal body", sl.Err(err))
		return err
	}

	accountId, err := uuid.Parse(event.AccountId)
	if err != nil {
		slog.Error("failed to parse account id", sl.Err(err))
		return err
	}

	serviceId, err := uuid.Parse(event.ServiceId)
	if err != nil {
		slog.Error("failed to parse service id", sl.Err(err))
		return err
	}

	if _, err := c.uc.CreateCabinet(ctx, &dto.NewCabinet{
		AccountId: accountId,
		ServiceId: serviceId,
	}); err != nil {
		slog.Error("failed to create cabinet", sl.Err(err))
		return err
	}

	return nil
}
