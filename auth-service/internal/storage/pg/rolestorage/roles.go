package rolestorage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/storage/pg"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

func (s *RoleStorage) Roles(ctx context.Context, userId uuid.UUID) ([]models.Role, error) {

	log := slog.With(sl.Method("credentialStorage.Roles"))

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
