package housestorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/mu/housing-service/internal/dto"
	"github.com/tehrelt/mu/housing-service/internal/storage/pg"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
)

func (s *HouseStorage) Create(ctx context.Context, in *dto.CreateHouse) (id uuid.UUID, err error) {

	log := slog.With(sl.Method("housestorage.Create"))

	log.Debug("creating house", slog.Any("create house dto", in))

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
	}

	defer func() {
		if err != nil {
			log.Warn("rollback transaction")
			_ = tx.Rollback()
			return
		}

		log.Debug("commiting transaction")
		err = tx.Commit()
	}()

	sql, args, err := sq.Insert(pg.HOUSE_TABLE).
		Columns("address", "rooms_qty", "residents_qty").
		Values(in.Address, in.RoomsQty, in.ResidentQty).
		Suffix("RETURNING (id)").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return
	}

	log.Debug("executing query", slog.String("sql", sql), slog.Any("args", args))

	var rawId string

	if err := tx.QueryRowContext(ctx, sql, args...).Scan(&rawId); err != nil {
		log.Error("failed to execute query", sl.Err(err))
	}
	log.Debug("user created, creating personal data entry", slog.String("id", rawId))
	id, err = uuid.Parse(rawId)
	if err != nil {
		log.Error("failed to parse id", slog.String("rawId", rawId), sl.Err(err))
		return
	}

	return
}
