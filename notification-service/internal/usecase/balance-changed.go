package usecase

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/notification-service/internal/events"
	"github.com/tehrelt/mu/notification-service/pkg/pb/accountpb"
	"go.opentelemetry.io/otel"
)

func (uc *UseCase) HandleBalanceChangedEvent(ctx context.Context, incoming *events.IncomingBalanceChanged) (err error) {

	fn := "HandleBalanceChangedEvent"
	log := uc.logger.With(sl.Method(fn))

	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	acc, err := uc.accountapi.Find(ctx, &accountpb.FindRequest{
		Id: incoming.AccountId,
	})
	if err != nil {
		log.Error("failed to find account", slog.String("accId", incoming.AccountId), sl.Err(err))
		return err
	}

	userid, _ := uuid.Parse(acc.Account.UserId)

	event := &events.BalanceChanged{
		EventHeader: events.NewEventHeader(events.EventBalanceChanged, userid.String()),
		NewBalance:  incoming.NewBalance,
		OldBalance:  incoming.OldBalance,
		Reason:      incoming.Reason,
		Address:     acc.Account.House.Address,
	}

	if err := uc.send(ctx, event); err != nil {
		log.Error("failed to send event", sl.Err(err))
		return err
	}

	return nil
}
