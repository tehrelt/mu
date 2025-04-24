package housestorage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu/housing-service/internal/models"
	"github.com/tehrelt/mu/housing-service/internal/storage/pg"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
)

func (s *HouseStorage) Find(ctx context.Context, houseId uuid.UUID) (*models.House, error) {

	log := slog.With(sl.Method("housestorage.Create"))

	log.Debug("find house", slog.Any("house_id", houseId))

	query, args, err := sq.Select("address", "rooms_qty", "residents_qty", "created_at", "updated_at").
		From(pg.HOUSE_TABLE).
		Where(sq.Eq{"id": houseId.String()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", slog.String("sql", query), slog.Any("args", args))

	house := new(models.House)
	if err := s.db.
		QueryRowContext(ctx, query, args...).
		Scan(&house.Address, &house.RoomsQuantity, &house.ResidentsQuantity, &house.CreatedAt, &house.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			log.Warn("no house found", sl.UUID("house_id", houseId))
		}
		log.Error("failed to execute query", sl.Err(err))
	}
	house.Id = houseId

	query, args, err = sq.Select("service_id").From(pg.CONNECTED_SERVICES_TABLE).Where(sq.Eq{"house_id": houseId.String()}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Error("failed to build query, to get services", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", slog.String("sql", query), slog.Any("args", args))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error("failed to execute query, to get services", sl.Err(err))
		return nil, err
	}

	for rows.Next() {
		var serviceId string
		if err := rows.Scan(&serviceId); err != nil {
			log.Error("failed to scan service_id", sl.Err(err))
			return nil, err
		}

		sid, err := uuid.Parse(serviceId)
		if err != nil {
			log.Error("failed to parse uuid", sl.Err(err))
			return nil, err
		}

		house.ConnectedServices = append(house.ConnectedServices, sid)
	}

	return house, nil
}
