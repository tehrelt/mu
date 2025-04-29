package amqp

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp091 "github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/housing-service/internal/dto"
	"github.com/tehrelt/mu/housing-service/internal/models"
	ratepb "github.com/tehrelt/mu/housing-service/pkg/pb/ratespb"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *AmqpConsumer) ConsumeConnectServiceEvent(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := c.manager.Consume(ctx, c.cfg.ConnectServiceQueue.Routing, c.handleConnectServiceEvent); err != nil {
				slog.Error("failed to consume", sl.Err(err))
				return err
			}
		}
	}
}

func (c *AmqpConsumer) handleConnectServiceEvent(ctx context.Context, msg amqp091.Delivery) (err error) {

	defer func() {
		if err != nil {
			msg.Nack(false, true)
			return
		}

		err = msg.Ack(false)
	}()

	body := msg.Body
	unmarshaled := dto.EventServiceConnect{}
	if err := json.Unmarshal(body, &unmarshaled); err != nil {
		slog.Error("failed to unmarshal body", sl.Err(err))
		return err
	}

	event, err := models.ParseEventConnectService(unmarshaled.HouseId, unmarshaled.ServiceId)
	if err != nil {
		slog.Error("failed to parse uuids", sl.Err(err))
		return err
	}

	if _, err := c.rateapi.Find(ctx, &ratepb.FindRequest{
		Id: event.ServiceId.String(),
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				slog.Error("service not found")
				return err
			}
		}

		slog.Error("failed to find service", sl.Err(err))
		return err
	}

	if err := c.storage.ConnectService(ctx, &dto.ConnectService{
		HouseId:   event.HouseId,
		ServiceId: event.ServiceId,
	}); err != nil {
		slog.Error("failed to connecting service", sl.Err(err))
		return err
	}

	if err := c.broker.PublishServiceConnectedEvent(ctx, &dto.EventServiceConnected{
		HouseId:   event.HouseId.String(),
		ServiceId: event.ServiceId.String(),
	}); err != nil {
		slog.Error("failed to publish service connected event", sl.Err(err))
		return err
	}

	return nil
}
