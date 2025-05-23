package integrationstorage

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/models"
	"github.com/tehrelt/mu/notification-service/internal/storage"
	"github.com/tehrelt/mu/notification-service/internal/storage/pg"
)

func (s *Storage) Create(ctx context.Context, integration *models.Integration) error {

	fn := "Create"
	log := s.logger.With(sl.Method(fn))

	query, args, err := sq.
		Insert(pg.TableIntegrations).
		Columns("user_id", "telegram_chat_id").
		Values(integration.UserId, integration.TelegramChatId).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed build query", sl.Err(err))
		return err
	}

	if _, err := s.pool.Exec(ctx, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return storage.ErrUserAlreadyExists
			}
		}
		log.Error("failed execute query", sl.Err(err))
		return err
	}

	return nil
}
