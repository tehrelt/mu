package accountstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/internal/models"
	"github.com/tehrelt/mu/account-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *AccountStorage) List(ctx context.Context, filters *dto.AccountFilters) (<-chan models.Account, error) {

	fn := "accountstorage.List"
	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(
		ctx,
		fn,
	)
	log := slog.With(sl.Method(fn))

	log.Debug("list payments", slog.Any("filters", filters))

	builder := sq.
		Select("id, user_id, house_id, balance, created_at, updated_at").
		From(pg.ACCOUNTS_TABLE)

	if filters.UserId != uuid.Nil {
		builder = builder.Where(sq.Eq{"user_id": filters.UserId})
		span.SetAttributes(attribute.String("user_id", filters.UserId.String()))
	}

	query, args, err := builder.
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	out := make(chan models.Account)
	errchan := make(chan error, 1)
	go func() {
		defer span.End()
		defer close(errchan)
		defer close(out)

		rows, err := s.db.QueryxContext(ctx, query, args...)
		if err != nil {
			log.Error("failed to execute query", sl.Err(err))
			span.RecordError(err)
			errchan <- err
			return
		}

		for rows.Next() {
			var acc models.Account
			if err := rows.Scan(
				&acc.Id,
				&acc.UserId,
				&acc.HouseId,
				&acc.Balance,
				&acc.CreatedAt,
				&acc.UpdatedAt,
			); err != nil {
				log.Error("failed to execute query", sl.Err(err))
				span.RecordError(err)
				errchan <- err
				continue
			}

			log.Debug("account sent to channel", slog.String("account_id", acc.Id))
			out <- acc
		}
	}()

	go func() {
		for range errchan {
			log.Error("failed to execute query", sl.Err(err))
		}
	}()

	return out, nil
}
