package servicestorage

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu/rate-service/internal/models"
	"github.com/tehrelt/mu/rate-service/internal/storage"
	"github.com/tehrelt/mu/rate-service/internal/storage/pg"
	"github.com/tehrelt/mu/rate-service/pkg/sl"
)

func (s *ServiceStorage) Update(ctx context.Context, in *models.UpdateServiceRate) error {
	log := slog.With(sl.Method("servicestorage.Update"))

	log.Debug("creating service")

	query, args, err := squirrel.Update(pg.SERVICES_TABLE).
		Set("rate", in.Rate).
		Where(squirrel.Eq{"id": in.Id.String()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Error("failed to get rows affected", sl.Err(err))
		return err
	}

	if affected == 0 {
		return storage.ErrServiceNotFound
	}

	return nil
}
