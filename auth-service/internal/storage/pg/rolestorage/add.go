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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return err
	}

	defer func() {
		if err != nil {
			log.Debug("rollback tx")
			_ = tx.Rollback()
			return
		}

		log.Debug("commit tx")
		err = tx.Commit()
		if err != nil {
			log.Warn("failed to commit transaction", sl.Err(err))
		}
	}()

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
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return err
	}

	return nil
}
