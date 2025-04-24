package accountstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *AccountStorage) Create(ctx context.Context, in *dto.CreateAccount) (id uuid.UUID, err error) {

	fn := "accountstorage.Create"
	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	defer span.End()
	span.SetAttributes(
		attribute.String("user_id", in.UserId.String()),
		attribute.String("house_id", in.HouseId.String()),
	)
	log := slog.With(sl.Method(fn))
	log.Debug("creating account", slog.Any("create account dto", in))

	query, args, err := sq.Insert(pg.ACCOUNTS_TABLE).
		Columns("user_id", "house_id").
		Values(in.UserId, in.HouseId).
		Suffix("RETURNING (id)").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return uuid.Nil, err
	}

	log.Debug("executing query", slog.String("sql", query), slog.Any("args", args))

	if err := s.db.QueryRowxContext(ctx, query, args...).Scan(&id); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return uuid.Nil, err
	}

	return id, nil
}
