package consumptionstorage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/internal/models"
	"github.com/tehrelt/mu/consumption-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (s *Storage) Logs(ctx context.Context, filters *dto.LogsFilters) ([]*models.ConsumptionLog, error) {

	fn := "ListLogs"
	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	log := s.logger.With(sl.Method(fn))

	len := 100

	builder := sq.
		Select("l.id", "l.amount", "l.payment_id", "l.cabinet_id", "l.created_at", "c.account_id", "c.service_id").
		From(fmt.Sprintf("%s as l", pg.ConsumptionLogTable)).
		Join(fmt.Sprintf("%s as c ON l.cabinet_id = c.id", pg.CabinetTable)).
		PlaceholderFormat(sq.Dollar).
		OrderBy("c.created_at DESC").
		Limit(100)

	if filters != nil {
		if !filters.Pagination.Nil() {
			if filters.Pagination.Limit > 0 {
				builder = builder.Limit(filters.Pagination.Limit)
				len = int(filters.Pagination.Limit)
			}
			if filters.Pagination.Offset > 0 {
				builder = builder.Offset(filters.Pagination.Offset)
			}
		}

		if filters.CabinetId != uuid.Nil {
			builder = builder.Where(sq.Eq{"l.cabinet_id": filters.CabinetId})
		}

		if filters.AccountId != uuid.Nil {
			builder = builder.Where(sq.Eq{"c.account_id": filters.AccountId})
		}

		if filters.ServiceId != uuid.Nil {
			builder = builder.Where(sq.Eq{"c.service_id": filters.ServiceId})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	span.SetAttributes(attribute.String("query", query))
	log.Debug("executing query", sl.Query(query), sl.Args(args))

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Error("failed to execute query", sl.Err(err), slog.Any("pgerr", pgErr))
		} else {
			log.Error("failed to execute query", sl.Err(err))
		}
		return nil, err
	}
	defer rows.Close()

	logs := make([]*models.ConsumptionLog, 0, len)

	for rows.Next() {
		record := &models.ConsumptionLog{}
		if err := rows.Scan(
			&record.Id,
			&record.Consumed,
			&record.PaymentId,
			&record.CabinetId,
			&record.CreatedAt,
			&record.AccountId,
			&record.ServiceId,
		); err != nil {
			log.Error("failed to scan row", sl.Err(err))
			return nil, err
		}

		logs = append(logs, record)
	}

	return logs, nil
}

func (s *Storage) CountLogs(ctx context.Context, filters *dto.LogsFilters) (uint64, error) {

	fn := "ListLogs"
	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	log := s.logger.With(sl.Method(fn))

	builder := sq.
		Select("count(*)").
		From(pg.ConsumptionLogTable).
		PlaceholderFormat(sq.Dollar)

	if filters != nil {
		if filters.CabinetId != uuid.Nil {
			builder = builder.Where(sq.Eq{"cabinet_id": filters.CabinetId})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return 0, err
	}

	span.SetAttributes(attribute.String("query", query))
	log.Debug("executing query", sl.Query(query), sl.Args(args))

	var total uint64
	if err := s.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return 0, err
	}

	return total, nil
}
