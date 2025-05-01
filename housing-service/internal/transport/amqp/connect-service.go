package amqp

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu/housing-service/internal/dto"
	"github.com/tehrelt/mu/housing-service/internal/models"
	ratepb "github.com/tehrelt/mu/housing-service/pkg/pb/ratespb"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *AmqpConsumer) handleConnectServiceEvent(ctx context.Context, msg *rmqmanager.TracedDelivery) (err error) {

	unmarshaled := dto.EventServiceConnect{}
	defer func() {
		if err != nil {
			slog.Info("reject event", slog.Any("event", unmarshaled))
			msg.Reject(false)
			return
		}

		err = msg.Ack(false)
	}()

	body := msg.Body
	if err := json.Unmarshal(body, &unmarshaled); err != nil {
		slog.Error("failed to unmarshal body", sl.Err(err))
		return err
	}

	event, err := models.ParseEventConnectService(unmarshaled.HouseId, unmarshaled.ServiceId, unmarshaled.AccountId)
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
		AccountId: event.AccountId.String(),
		ServiceId: event.ServiceId.String(),
	}); err != nil {
		slog.Error("failed to publish service connected event", sl.Err(err))
		return err
	}

	return nil
}
