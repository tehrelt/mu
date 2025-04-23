package amqp

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/pkg/pb/billingpb"

	"github.com/tehrelt/mu-lib/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errPaymentNotFound = errors.New("invalid payment id")

func (c *AmqpConsumer) handlePaymentStatusChangedEvent(ctx context.Context, msg amqp091.Delivery) (err error) {

	defer func() {
		if err != nil && !errors.Is(err, errPaymentNotFound) {
			msg.Nack(false, true)
			return
		}

		err = msg.Ack(false)
	}()

	body := msg.Body
	unmarshaled := struct {
		AccountId string            `json:"accountId"`
		PaymentId string            `json:"paymentId"`
		NewStatus dto.PaymentStatus `json:"newStatus"`
	}{}

	if err := json.Unmarshal(body, &unmarshaled); err != nil {
		slog.Error("failed to unmarshal body", sl.Err(err))
		return err
	}

	accId, err := uuid.Parse(unmarshaled.AccountId)
	if err != nil {
		slog.Error("failed to parse account id", sl.Err(err))
		return err
	}

	payId, err := uuid.Parse(unmarshaled.PaymentId)
	if err != nil {
		slog.Error("failed to parse payment id", sl.Err(err))
		return err
	}

	event := &dto.EventPaymentStatusChanged{
		AccountId: accId,
		PaymentId: payId,
		NewStatus: unmarshaled.NewStatus,
	}

	if event.NewStatus != dto.PaymentStatusPaid {
		slog.Info("incoming payment status change is not paid, skipping...")
		return
	}

	acc, err := c.storage.Find(ctx, event.AccountId)
	if err != nil {
		slog.Error("failed to find account", sl.Err(err))
		return err
	}

	bill, err := c.billingApi.Find(ctx, &billingpb.FindRequest{Id: event.PaymentId.String()})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				slog.Error("payment not found, skipping...")
				return errPaymentNotFound
			}
		}
		slog.Error("failed to find payment", sl.Err(err))
		return err
	}

	oldBalance := acc.Balance
	acc.DeltaBalance(bill.Payment.Amount)

	if _, err := c.storage.Update(ctx, &dto.UpdateAccount{
		Id:         event.AccountId,
		NewBalance: acc.Balance,
	}); err != nil {
		slog.Error("failed to update account", sl.Err(err))
		return err
	}

	slog.Info("account balance updated", sl.UUID("account_id", accId))

	if err := c.broker.PublishBalanceChanged(ctx, &dto.EventBalanceChanged{
		AccountId:  acc.Id,
		NewBalance: acc.Balance,
		OldBalance: oldBalance,
		Reason:     "no reason",
	}); err != nil {
		slog.Error("failed to publish balance changed event")
		return err
	}

	return nil
}
