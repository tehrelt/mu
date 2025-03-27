package paymentstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu/billing-service/internal/dto"
	"github.com/tehrelt/mu/billing-service/internal/models"
	"github.com/tehrelt/mu/billing-service/internal/storage/pg"
	"github.com/tehrelt/mu/billing-service/pkg/sl"
)

func (s *PaymentStorage) List(ctx context.Context, filters *dto.PaymentFilters, out chan<- models.Payment) error {
	defer close(out)

	log := slog.With(sl.Method("paymentstorage.Create"))

	log.Debug("list payments", slog.Any("filters", filters))

	builder := sq.
		Select("*").
		From(pg.PAYMENTS_TABLE)

	if filters.AccountId != uuid.Nil {
		builder = builder.Where(sq.Eq{"account_id": filters.AccountId})
	}
	if filters.Status != models.PaymentStatusNil {
		builder = builder.Where(sq.Eq{"status": filters.Status})
	}
	if filters.AmountRange.Nil() {
		if filters.AmountRange.Max != 0 {
			builder = builder.Where(sq.Lt{"amount": filters.AmountRange.Max})
		}
		if filters.AmountRange.Min != 0 {
			builder = builder.Where(sq.Gt{"amount": filters.AmountRange.Min})
		}
	}

	query, args, err := builder.
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return err
	}

	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(
			&payment.Id,
			&payment.AccountId,
			&payment.Amount,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		); err != nil {
			log.Error("failed to execute query", sl.Err(err))
			return err
		}

		out <- payment
	}

	return nil
}
