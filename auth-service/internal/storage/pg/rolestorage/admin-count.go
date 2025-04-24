package rolestorage

import (
	"context"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/internal/storage/pg"
)

func (s *RoleStorage) Count(ctx context.Context, role ...models.Role) (int64, error) {

	fn := "rolestorage.Count"
	log := slog.With(slog.String("fn", fn))
	var count int64

	builder := sq.
		Select("COUNT(*)").
		From(pg.ROLES_TABLE).
		Where(sq.Eq{"role": role}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return count, err
	}

	log.Debug("exec", sl.Query(query), sl.Args(args))

	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		log.Error("failed to scan count", sl.Err(err))
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return count, nil
}
