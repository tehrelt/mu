package amqp

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/pkg/pb/housepb"
	"github.com/tehrelt/mu/account-service/pkg/pb/ticketpb"

	"github.com/tehrelt/mu-lib/sl"
)

var errTicketNotFound = errors.New("ticket not found")

func (c *AmqpConsumer) handleTicketStatusChanged(ctx context.Context, msg amqp091.Delivery) (err error) {

	defer func() {
		if err != nil && !errors.Is(err, errTicketNotFound) {
			msg.Nack(false, true)
			return
		}

		err = msg.Ack(false)
	}()

	fn := "consumer.handleTicketStatusChanged"
	log := slog.With(sl.Method(fn))

	body := msg.Body
	event := &dto.EventTicketStatusChanged{}
	if err := json.Unmarshal(body, &event); err != nil {
		slog.Error("failed to unmarshal body", sl.Err(err))
		return err
	}
	log.Debug("unmarshal body", slog.Any("event", event))

	if event.Status != "approved" {
		log.Info("event rejected", slog.Any("event", event))
		return nil
	}

	log.Debug("try to find ticket", slog.String("ticketId", event.TicketId))
	ticket, err := c.ticketApi.Find(ctx, &ticketpb.FindRequest{
		Id: event.TicketId,
	})
	if err != nil {
		slog.Error("failed to find ticket", sl.Err(err))
		return err
	}

	var accticket *ticketpb.Ticket_Account

	switch w := ticket.Ticket.Payload.(type) {
	case *ticketpb.Ticket_Account:
		log.Debug("ticket of new account type", slog.Any("payload", w))
		accticket = w
	default:
		slog.Error("unknown ticket type", sl.Err(err))
		return err
	}

	userId, err := uuid.Parse(ticket.Ticket.Header.CreatedBy)
	if err != nil {
		slog.Error("failed to parse user id", sl.Err(err))
		return err
	}

	newHouseReq := &housepb.CreateRequest{
		Address:      accticket.Account.HouseAdress,
		RoomsQty:     rand.Int63n(5),
		ResidentsQty: rand.Int63n(5),
	}
	log.Debug("creating house", slog.Any("key string", newHouseReq))
	house, err := c.houseApi.Create(ctx, newHouseReq)
	if err != nil {
		slog.Error("failed to create house", sl.Err(err))
		return err
	}
	log.Debug("house created", slog.Any("house", house))

	houseId, err := uuid.Parse(house.Id)
	if err != nil {
		slog.Error("failed to parse house id", sl.Err(err))
		return err
	}

	log.Debug("creating account", slog.Any("userId", userId), slog.Any("houseId", houseId))
	accid, err := c.storage.Create(ctx, &dto.CreateAccount{
		UserId:  userId,
		HouseId: houseId,
	})
	if err != nil {
		slog.Error("failed to create account", sl.Err(err))
		return err
	}
	log.Debug("account created", slog.Any("accountId", accid))

	return nil
}
