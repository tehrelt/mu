package servicestorage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu/rate-service/internal/models"
	"github.com/tehrelt/mu/rate-service/internal/storage"
	"github.com/tehrelt/mu/rate-service/internal/storage/pg"
	"github.com/tehrelt/mu/rate-service/pkg/sl"
)

func (s *ServiceStorage) Find(ctx context.Context, id uuid.UUID) (*models.Service, error) {
	log := slog.With(sl.Method("servicestorage.Find"))

	log.Debug("searching for service", sl.UUID("service_id", id))

	query, args, err := squirrel.Select("*").
		From(pg.SERVICES_TABLE).
		Where(squirrel.Eq{"id": id.String()}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	service := new(models.Service)
	if err := s.db.GetContext(ctx, service, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("service not found", sl.UUID("service_id", id))
			return nil, storage.ErrServiceNotFound
		}
		log.Error("failed to execute query", sl.Err(err))
		return nil, err
	}

	return service, nil
}
