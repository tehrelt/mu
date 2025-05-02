package billingapi

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/billingpb"
)

func (a *Api) Charge(ctx context.Context, bill *dto.Charge) (uuid.UUID, error) {
	fn := "Charge"
	logger := a.logger.With(sl.Method(fn))

	req := &billingpb.CreateRequest{
		AccountId: bill.AccountId.String(),
		Amount:    -int64(bill.Amount),
		Message:   bill.Message,
	}

	logger.Debug("creating payment", slog.Any("req", req))
	res, err := a.client.Create(ctx, req)
	if err != nil {
		logger.Error("failed create payment", sl.Err(err))
		return uuid.Nil, err
	}
	logger.Debug("payment created", slog.Any("payment", res))

	logger.Debug("pay that bill", slog.Any("bill", res))
	if _, err := a.client.Pay(ctx, &billingpb.PayRequest{
		PaymentId: res.Id,
	}); err != nil {
		logger.Error("failed to pay", sl.Err(err))
		return uuid.Nil, err
	}

	return uuid.MustParse(res.Id), nil
}
