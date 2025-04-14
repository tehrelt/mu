package credentialstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/auth-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
)

func (s *CredentialStorage) Password(ctx context.Context, userId uuid.UUID) (string, error) {

	fn := "credentialstorage.Password"
	log := slog.With(sl.Method(fn))

	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, fn)
	defer span.End()

	log.Debug("get credentials", slog.String("userId", userId.String()))

	query, args, err := sq.Select("hashed_password").From(pg.CREDENTIALS_TABLE).Where(sq.Eq{"id": userId.String()}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return "", err
	}

	var hashPassword string
	row := s.db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(&hashPassword); err != nil {
		log.Error("failed to get credentials", sl.Err(err))
		return "", err
	}

	return hashPassword, nil
}
