package rolestorage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
)

func (s *RoleStorage) Roles(ctx context.Context, userId uuid.UUID) ([]models.Role, error) {

	fn := "rolestorage.Roles"
	log := slog.With(sl.Method(fn))

	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, fn)
	defer span.End()

	query, args, err := sq.Select("role").
		From(pg.ROLES_TABLE).
		Where(sq.Eq{"user_id": userId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	roles := make([]models.Role, 0, 2)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role string
		if err = rows.Scan(&role); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return roles, nil
			}

			log.Error("failed to scan row", sl.Err(err))
			return nil, err
		}

		roles = append(roles, models.Role(role))
	}

	return roles, nil
}
