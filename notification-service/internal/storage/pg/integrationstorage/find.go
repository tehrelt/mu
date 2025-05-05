package integrationstorage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/models"
	"github.com/tehrelt/mu/notification-service/internal/storage/pg"
)

func (s *Storage) Find(ctx context.Context, userId uuid.UUID) (*models.Integration, error) {

	fn := "Find"
	log := s.logger.With(sl.Method(fn))

	query, args, err := sq.
		Select("user_id", "telegram_chat_id").
		From(pg.TableIntegrations).
		Where(sq.Eq{"user_id": userId.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing", sl.Query(query), sl.Args(args))
	out := new(models.Integration)
	row := s.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&out.UserId, &out.TelegramChatId); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		log.Error("failed scan row", sl.Err(err))
		return nil, err
	}

	return out, nil
}
