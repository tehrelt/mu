package paymentstorage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/billing-service/internal/models"
	"github.com/tehrelt/mu/billing-service/internal/storage"
	"github.com/tehrelt/mu/billing-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
)

func (s *PaymentStorage) Find(ctx context.Context, paymentId uuid.UUID) (*models.Payment, error) {

	fn := "paymentstorage.Find"
	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, fn)
	defer span.End()
	log := slog.With(sl.Method(fn))

	log.Debug("find payment", slog.Any("payment_id", paymentId))

	query, args, err := sq.
		Select("id", "account_id", "amount", "message", "status", "created_at", "updated_at").
		From(pg.PAYMENTS_TABLE).
		Where(sq.Eq{"id": paymentId.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	payment := new(models.Payment)
	if err := s.db.
		QueryRowContext(ctx, query, args...).
		Scan(
			&payment.Id,
			&payment.AccountId,
			&payment.Amount,
			&payment.Message,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("no payments found", sl.UUID("payment_id", paymentId))
			return nil, storage.ErrPaymentNotFound
		}

		log.Error("failed to execute query", sl.Err(err))
		return nil, err
	}

	return payment, nil
}
