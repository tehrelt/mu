package consumptionstorage

import (
	"context"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/internal/models"
	"github.com/tehrelt/mu/consumption-service/internal/storage"
	"github.com/tehrelt/mu/consumption-service/internal/storage/pg"
	"github.com/tehrelt/mu/consumption-service/internal/usecase"
)

var _ usecase.ConsumptionStorage = (*Storage)(nil)

type Storage struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// Create implements usecase.ConsumptionStorage.
func (s *Storage) Create(ctx context.Context, in *dto.NewCabinet) (*models.Cabinet, error) {
	fn := "Create"
	log := s.logger.With(sl.Method(fn))

	query, args, err := sq.
		Insert(pg.CabinetTable).
		Columns("account_id", "service_id").
		Values(in.AccountId, in.ServiceId).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed build query", sl.Err(err))
		return nil, err
	}

	row := s.pool.QueryRow(ctx, query, args...)
	var out models.Cabinet
	err = row.Scan(&out.Id, &out.AccountId, &out.ServiceId, &out.Consumed, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		log.Error("failed scan row", sl.Err(err))
		return nil, err
	}

	return &out, nil
}

// Find implements usecase.ConsumptionStorage.
func (s *Storage) Find(ctx context.Context, criteria *dto.FindCabinet) (*models.Cabinet, error) {
	fn := "Find"
	log := s.logger.With(sl.Method(fn))

	builder := sq.
		Select("*").
		From(pg.CabinetTable).
		PlaceholderFormat(sq.Dollar)

	if criteria.Id != uuid.Nil {
		builder = builder.Where(sq.Eq{"id": criteria.Id})
	} else if criteria.AccountId != uuid.Nil && criteria.ServiceId != uuid.Nil {
		builder = builder.Where(sq.Eq{
			"account_id": criteria.AccountId,
			"service_id": criteria.ServiceId,
		})
	} else {
		return nil, storage.ErrInvalidCriteria
	}
	query, args, err := builder.
		ToSql()
	if err != nil {
		log.Error("failed build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))
	row := s.pool.QueryRow(ctx, query, args...)
	var out models.Cabinet
	err = row.Scan(&out.Id, &out.AccountId, &out.ServiceId, &out.Consumed, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		log.Error("failed scan row", sl.Err(err))
		return nil, err
	}

	return &out, nil
}

// Log implements usecase.ConsumptionStorage.
func (s *Storage) Log(ctx context.Context, in *dto.NewConsumeLog) (*models.ConsumptionLog, error) {
	fn := "Log"
	log := s.logger.With(sl.Method(fn))

	query, args, err := sq.
		Insert(pg.ConsumptionLogTable).
		Columns("cabinet_id", "amount", "payment_id").
		Values(in.CabinetId, in.Amount, in.PaymentId).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))
	row := s.pool.QueryRow(ctx, query, args...)
	var out models.ConsumptionLog
	err = row.Scan(&out.Id, &out.Amount, &out.PaymentId, &out.CabinetId, &out.CreatedAt)
	if err != nil {
		log.Error("failed scan row", sl.Err(err))
		return nil, err
	}

	return &out, nil
}

// Update implements usecase.ConsumptionStorage.
func (s *Storage) Update(ctx context.Context, in *dto.UpdateCabinet) (*models.Cabinet, error) {
	fn := "Update"
	log := s.logger.With(sl.Method(fn))

	builder := sq.
		Update(pg.CabinetTable).
		Where(sq.Eq{"id": in.Id}).
		Set("updated_at", time.Now()).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	if in.AmountDelta != 0 {
		builder = builder.Set("consumed", sq.Expr("consumed + ?", in.AmountDelta))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))
	row := s.pool.QueryRow(ctx, query, args...)
	var out models.Cabinet
	err = row.Scan(&out.Id, &out.AccountId, &out.ServiceId, &out.Consumed, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		log.Error("failed scan row", sl.Err(err))
		return nil, err
	}

	return &out, nil
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool:   pool,
		logger: slog.With(sl.Module("consumptionstorage")),
	}
}
