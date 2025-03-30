package sessionstorage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/config"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/lib/tracer"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/services/authservice"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/storage"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
	"go.opentelemetry.io/otel"

	"github.com/redis/go-redis/v9"
)

var _ authservice.SessionsStorage = (*SessionsStorage)(nil)

const traceKey = "redis:session-storage"

type SessionsStorage struct {
	db     *redis.Client
	cfg    *config.Config
	logger *slog.Logger
}

func (s *SessionsStorage) Save(ctx context.Context, userId uuid.UUID, token string) error {

	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, traceKey)
	defer span.End()

	fn := "redis.Save"
	log := s.logger.With(slog.String("userId", userId.String()), sl.Method(fn))

	log.Debug("Saving session")

	stat := s.db.Set(ctx, userId.String(), token, time.Duration(s.cfg.Jwt.RefreshTTL)*time.Minute)
	if err := stat.Err(); err != nil {
		log.Error("cannot save session", sl.Err(err))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *SessionsStorage) Check(ctx context.Context, userId uuid.UUID, refreshToken string) error {

	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, traceKey)
	defer span.End()

	fn := "redis.Check"
	log := s.logger.With(slog.String("userId", userId.String()), sl.Method(fn))

	log.Debug("Checking session")

	stat := s.db.Get(ctx, userId.String())
	if err := stat.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return fmt.Errorf("%s: %w", fn, storage.ErrSessionNotFound)
		}

		log.Error("cannot check session", sl.Err(err))
		return fmt.Errorf("%s: %w", fn, err)
	}

	if stat.Val() != refreshToken {
		log.Error("invalid session", slog.String("user_id", userId.String()), slog.String("refresh_token", refreshToken), slog.String("session_token", stat.Val()))
		return fmt.Errorf("%s: %w", fn, storage.ErrSessionInvalid)
	}

	return nil
}

func (s *SessionsStorage) Delete(ctx context.Context, userId uuid.UUID) error {

	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, traceKey)
	defer span.End()

	fn := "redis.Delete"
	log := s.logger.With(sl.Method(fn), slog.String("userId", userId.String()))

	log.Debug("Deleting session")

	stat := s.db.Del(ctx, userId.String())
	if err := stat.Err(); err != nil {
		log.Error("cannot delete session", sl.Err(err))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func New(db *redis.Client, cfg *config.Config) *SessionsStorage {
	return &SessionsStorage{
		db:     db,
		cfg:    cfg,
		logger: slog.With(sl.Module("redis.SessionStorage")),
	}
}
