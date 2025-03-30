package accountstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/internal/models"
	"github.com/tehrelt/mu/account-service/internal/storage/pg"
	"github.com/tehrelt/mu/account-service/pkg/sl"
)

func (s *AccountStorage) List(ctx context.Context, filters *dto.AccountFilters, out chan<- models.Account) error {
	defer close(out)

	log := slog.With(sl.Method("paymentstorage.Create"))

	log.Debug("list payments", slog.Any("filters", filters))

	builder := sq.
		Select("id, user_id, house_id, balance, created_at, updated_at").
		From(pg.ACCOUNTS_TABLE)

	query, args, err := builder.
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return err
	}

	for rows.Next() {
		var acc models.Account
		if err := rows.Scan(
			&acc.Id,
			&acc.UserId,
			&acc.HouseId,
			&acc.Balance,
			&acc.CreatedAt,
			&acc.UpdatedAt,
		); err != nil {
			log.Error("failed to execute query", sl.Err(err))
			return err
		}

		out <- acc
	}

	return nil
}
