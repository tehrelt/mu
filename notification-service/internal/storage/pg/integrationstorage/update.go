package integrationstorage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/models"
	"github.com/tehrelt/mu/notification-service/internal/storage/pg"
)

func (s *Storage) Update(ctx context.Context, integration *models.Integration) error {

	fn := "Update"
	log := s.logger.With(sl.Method(fn))

	query, args, err := sq.
		Update(pg.TableIntegrations).
		Set("telegram_chat_id", integration.TelegramChatId).
		Where(sq.Eq{"user_id": integration.UserId.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed build query", sl.Err(err))
		return err
	}

	if _, err := s.pool.Exec(ctx, query, args...); err != nil {
		log.Error("failed execute query", sl.Err(err))
		return err
	}

	return nil
}
