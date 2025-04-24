package servicestorage

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu/rate-service/internal/models"
	"github.com/tehrelt/mu/rate-service/internal/storage/pg"
	"github.com/tehrelt/mu/rate-service/pkg/sl"
)

func (s *ServiceStorage) List(ctx context.Context, out chan<- *models.Service) error {
	log := slog.With(sl.Method("servicestorage.List"))

	log.Debug("listing all services")

	query, args, err := squirrel.Select("*").
		From(pg.SERVICES_TABLE).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to query", sl.Err(err))
		return err
	}
	for rows.Next() {
		service := new(models.Service)
		err = rows.Scan(&service.Id, &service.Name, &service.MeasureUnit, &service.Rate, &service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			log.Error("failed to scan row", sl.Err(err))
			return err
		}

		out <- service
	}
	close(out)

	return nil
}
