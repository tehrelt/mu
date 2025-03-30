package accountstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/internal/storage/pg"
	"github.com/tehrelt/mu/account-service/pkg/sl"
)

func (s *AccountStorage) Create(ctx context.Context, in *dto.CreateAccount) (id uuid.UUID, err error) {

	log := slog.With(sl.Method("accountstorage.Create"))

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
