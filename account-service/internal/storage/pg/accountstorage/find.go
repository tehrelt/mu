package accountstorage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/account-service/internal/models"
	"github.com/tehrelt/mu/account-service/internal/storage"
	"github.com/tehrelt/mu/account-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *AccountStorage) Find(ctx context.Context, id uuid.UUID) (*models.Account, error) {

	fn := "accountstorage.Find"
	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(
		ctx,
		fn,
		trace.WithAttributes(
			attribute.String("account_id", id.String()),
		),
	)
	defer span.End()
	log := slog.With(sl.Method(fn))

	log.Debug("find account", slog.Any("account_id", id))

	query, args, err := sq.
		Select("id, user_id, house_id, balance, created_at, updated_at").
		From(pg.ACCOUNTS_TABLE).
		Where(sq.Eq{"id": id.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	acc := new(models.Account)
	if err := s.db.
		QueryRowxContext(ctx, query, args...).
		Scan(
			&acc.Id,
			&acc.UserId,
			&acc.HouseId,
			&acc.Balance,
			&acc.CreatedAt,
			&acc.UpdatedAt,
		); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("no account found", sl.UUID("payment_id", id))
			return nil, storage.ErrAccountNotFound
		}

		log.Error("failed to execute query", sl.Err(err))
		return nil, err
	}

	return acc, nil
}
