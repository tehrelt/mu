package servicestorage

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu/rate-service/internal/models"
	"github.com/tehrelt/mu/rate-service/internal/storage/pg"
	"github.com/tehrelt/mu/rate-service/pkg/sl"
)

func (s *ServiceStorage) List(ctx context.Context, f *models.RateFilters) (<-chan *models.Service, error) {
	log := slog.With(sl.Method("servicestorage.List"))

	log.Debug("listing all services")

	builder := squirrel.Select("id, s_name, measure_unit, rate, created_at, updated_at, s_type").
		From(pg.SERVICES_TABLE).
		PlaceholderFormat(squirrel.Dollar)

	if f != nil {
		if f.Type != nil {
			builder = builder.Where(squirrel.Eq{"s_type": f.Type})
		}
	}

	query, args, err := builder.
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to query", sl.Err(err))
		return nil, err
	}

	out := make(chan *models.Service, 10)
	go func() {
		defer rows.Close()
		defer close(out)

		for rows.Next() {
			service := new(models.Service)
			err = rows.Scan(&service.Id, &service.Name, &service.MeasureUnit, &service.Rate, &service.CreatedAt, &service.UpdatedAt, &service.Type)
			if err != nil {
				log.Error("failed to scan row", sl.Err(err))
				continue
			}

			out <- service
		}
	}()

	return out, nil
}
