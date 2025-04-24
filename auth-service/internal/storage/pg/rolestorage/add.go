package rolestorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/storage/pg"

	"go.opentelemetry.io/otel"
)

func (s *RoleStorage) Add(ctx context.Context, in *dto.UserRoles) (err error) {

	fn := "rolestorage.Add"
	log := slog.With(sl.Method(fn))

	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, fn)
	defer span.End()

	builder := sq.Insert(pg.ROLES_TABLE).Columns("user_id", "role").PlaceholderFormat(sq.Dollar)

	for _, role := range in.Roles {
		if !role.Valid() {
			log.Warn("invalid role in request", slog.String("role", string(role)))
			continue
		}

		builder = builder.Values(in.UserId, role)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Warn("failed to build query for saving roles")
		return err
	}

	log.Debug("saving roles", slog.String("query", query), slog.Any("args", args))
	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return err
	}

	return nil
}
