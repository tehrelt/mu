package amqp

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/pkg/pb/ticketpb"

	"github.com/tehrelt/mu-lib/rmqmanager"
	"github.com/tehrelt/mu-lib/sl"
)

func (c *AmqpConsumer) handleConnectServiceEvent(ctx context.Context, msg *rmqmanager.TracedDelivery) (err error) {

	defer func() {
		if err != nil && !errors.Is(err, errTicketNotFound) {
			msg.Nack(false, true)
			return
		}

		err = msg.Ack(false)
	}()

	fn := "consumer.handleTicketConnectServiceStatusChanged"
	log := slog.With(sl.Method(fn))

	body := msg.Body
	incomingEvent := &dto.EventTicketStatusChanged{}
	if err := json.Unmarshal(body, &incomingEvent); err != nil {
		slog.Error("failed to unmarshal body", sl.Err(err))
		return err
	}
	log.Debug("unmarshal body", slog.Any("event", incomingEvent))

	if incomingEvent.Status != "approved" {
		log.Info("event rejected", slog.Any("event", incomingEvent))
		return nil
	}

	log.Debug("try to find ticket", slog.String("ticketId", incomingEvent.TicketId))
	ticketResponse, err := c.ticketApi.Find(ctx, &ticketpb.FindRequest{
		Id: incomingEvent.TicketId,
	})
	if err != nil {
		slog.Error("failed to find ticket", sl.Err(err))
		return err
	}

	var ticket *ticketpb.Ticket_ConnectService

	switch w := ticketResponse.Ticket.Payload.(type) {
	case *ticketpb.Ticket_ConnectService:
		log.Debug("ticket of connect service type", slog.Any("payload", w))
		ticket = w
	default:
		log.Error("unknown ticket type", sl.Err(err))
		return err
	}

	accid, err := uuid.Parse(ticket.ConnectService.AccountId)
	if err != nil {
		log.Error("failed to parse account id", sl.Err(err))
		return err
	}

	account, err := c.storage.Find(ctx, accid)
	if err != nil {
		log.Error("failed to find account", sl.Err(err))
		return err
	}

	if err := c.broker.PublishConnectServiceRequest(ctx, &dto.EventServiceConnect{
		AccountId: account.Id,
		HouseId:   account.HouseId,
		ServiceId: ticket.ConnectService.ServiceId,
	}); err != nil {
		log.Error("failed to publish connect service request", sl.Err(err))
		return err
	}

	return nil
}
