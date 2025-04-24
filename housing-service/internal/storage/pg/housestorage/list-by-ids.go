package housestorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tehrelt/mu/housing-service/internal/models"
	"github.com/tehrelt/mu/housing-service/internal/storage/pg"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
)

func (s *HouseStorage) ListByIds(ctx context.Context, ids []uuid.UUID) ([]*models.House, error) {

	log := slog.With(sl.Method("housestorage.Create"))

	log.Debug("list houses by ids", slog.Any("ids", ids))

	builder := sq.Select("address", "rooms_qty", "residents_qty", "created_at", "updated_at").
		From(pg.HOUSE_TABLE).
		PlaceholderFormat(sq.Dollar)

	for _, id := range ids {
		builder = builder.Where(sq.Eq{"id": id.String()})
	}

	query, args, err := builder.ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", slog.String("sql", query), slog.Any("args", args))

	houses := make([]houseDbModel, 0)
	if err := s.db.SelectContext(ctx, houses, query, args...); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, err
	}

	return lo.Map(houses, func(item houseDbModel, _ int) *models.House {
		h, err := item.ToHouse()
		if err != nil {
			return nil
		}

		return h
	}), nil
}
