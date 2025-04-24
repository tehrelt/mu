package housestorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu/housing-service/internal/dto"
	"github.com/tehrelt/mu/housing-service/internal/storage/pg"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
)

func (s *HouseStorage) ConnectService(ctx context.Context, in *dto.ConnectService) error {
	fn := "housestorage.find"
	log := slog.With(sl.Method(fn))

	query, args, err := sq.
		Insert(pg.CONNECTED_SERVICES_TABLE).
		Columns("house_id", "service_id").
		Values(in.HouseId, in.ServiceId).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return err
	}

	log.Debug("query", sl.Query(query), sl.Args(args))

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		log.Error("failed to exec query result", sl.Err(err))
		return err
	}

	return nil
}
