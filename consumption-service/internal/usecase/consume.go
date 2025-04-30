package usecase

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
)

func (uc *UseCase) Consume(ctx context.Context, in *dto.NewConsume) (*dto.ConsumeCreated, error) {
	fn := "usecase.Consume"
	logger := uc.logger.With(sl.Method(fn))

	cabinet, err := uc.find(ctx, &in.FindCabinet)
	if err != nil {
		return nil, err
	}

	service, err := uc.findService(ctx, cabinet.ServiceId)
	if err != nil {
		return nil, err
	}

	amount := in.Consumed * service.Rate
	logger.Debug(
		"calculated amount",
		slog.Uint64("consumed", in.Consumed),
		slog.Uint64("rate", service.Rate),
		slog.Uint64("amount", amount),
	)

	charge := &dto.Charge{
		AccountId: cabinet.AccountId,
		ServiceId: service.Id,
		Amount:    amount,
	}

	logger.Debug("creating charge", slog.Any("charge", charge))
	paymentId, err := uc.billingProvider.Charge(ctx, charge)
	if err != nil {
		return nil, ErrPaymentServiceUnavailable
	}
	newLog := &dto.NewConsumeLog{
		Consumed:  in.Consumed,
		PaymentId: paymentId,
		CabinetId: cabinet.Id,
	}

	logger.Debug("creating log record", slog.Any("log", newLog))
	log, err := uc.storage.Log(ctx, newLog)
	if err != nil {
		return nil, err
	}

	updateIn := &dto.UpdateCabinet{
		Id:            cabinet.Id,
		ConsumedDelta: in.Consumed,
	}
	logger.Debug("updating cabinet", slog.Any("updateDto", updateIn))
	if _, err := uc.storage.Update(ctx, updateIn); err != nil {
		return nil, err
	}
	ret := &dto.ConsumeCreated{
		LogId:     log.Id,
		PaymentId: paymentId,
		Amount:    amount,
	}

	logger.Info("consume record created", slog.Any("ret", ret))
	return ret, nil
}
