package paymentstorage

import (
	"context"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu/billing-service/internal/dto"
	"github.com/tehrelt/mu/billing-service/internal/models"
	"github.com/tehrelt/mu/billing-service/internal/storage/pg"
	"github.com/tehrelt/mu/billing-service/pkg/sl"
	"github.com/tehrelt/mu/billing-service/tracer"
	"go.opentelemetry.io/otel"
)

func (s *PaymentStorage) Update(ctx context.Context, in *dto.UpdatePayment) (*models.Payment, error) {

	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, traceKey)
	defer span.End()

	log := slog.With(sl.Method("paymentstorage.Create"))

	log.Debug("creating payment", slog.Any("create house dto", in))

	oldpayment, err := s.Find(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	builder := sq.Update(pg.PAYMENTS_TABLE).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": in.Id}).
		Set("updated_at", time.Now())

	if in.Amount != nil {
		builder = builder.Set("amount", *in.Amount)
	}

	if in.Status != nil {
		builder = builder.Set("status", *in.Status)
	}

	query, args, err := builder.
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, err
	}

	return oldpayment, nil
}
