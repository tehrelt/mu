package servicestorage

import (
	"context"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu/rate-service/internal/models"
	"github.com/tehrelt/mu/rate-service/internal/storage/pg"
	"github.com/tehrelt/mu/rate-service/pkg/sl"
)

func (s *ServiceStorage) Create(ctx context.Context, in *models.CreateService) (uuid.UUID, error) {
	log := slog.With(sl.Method("servicestorage.Create"))

	log.Debug("creating service", slog.Any("in", in))

	query, args, err := squirrel.Insert(pg.SERVICES_TABLE).
		Columns("s_name", "measure_unit", "rate", "s_type").
		Values(in.Name, in.MeasureUnit, in.Rate, in.Type).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return uuid.Nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	var createdId string
	if err := s.db.QueryRowContext(ctx, query, args...).Scan(&createdId); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return uuid.Nil, err
	}

	id, err := uuid.Parse(createdId)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))
		return uuid.Nil, err
	}

	return id, nil
}
