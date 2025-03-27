package paymentstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu/billing-service/internal/dto"
	"github.com/tehrelt/mu/billing-service/internal/storage/pg"
	"github.com/tehrelt/mu/billing-service/pkg/sl"
)

func (s *PaymentStorage) Create(ctx context.Context, in *dto.CreatePayment) (id uuid.UUID, err error) {

	log := slog.With(sl.Method("paymentstorage.Create"))

	log.Debug("creating payment", slog.Any("create house dto", in))

	query, args, err := sq.Insert(pg.PAYMENTS_TABLE).
		Columns("account_id", "amount").
		Values(in.AccountId, in.Amount).
		Suffix("RETURNING (id)").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return uuid.Nil, err
	}

	log.Debug("executing query", slog.String("sql", query), slog.Any("args", args))

	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return uuid.Nil, err
	}

	return id, nil
}
