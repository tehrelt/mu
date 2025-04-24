package userstorage

import (
	"context"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/user-service/internal/models"
	"github.com/tehrelt/mu/user-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
)

func (s *UserStorage) List(ctx context.Context, filters *models.UserFilters) (<-chan *models.User, error) {
	fn := "userstorage.List"
	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, fn)
	log := slog.With(sl.Method(fn))

	out := make(chan *models.User, 10)

	builder := sq.Select("u.id, u.last_name, u.first_name, u.middle_name, u.email, pd.phone, pd.snils, pd.passport_number, pd.passport_series, u.created_at, u.updated_at").
		From(fmt.Sprintf("%s u", pg.USERS)).
		Join(fmt.Sprintf("%s pd ON pd.user_id = u.id", pg.PERSONAL_DATA)).
		PlaceholderFormat(sq.Dollar).
		Limit(filters.Limit).
		Offset(filters.Offset).
		OrderBy("u.created_at ASC")

	query, args, err := builder.
		ToSql()

	if err != nil {
		log.Warn("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("failed to query", sl.Err(err))
		return nil, err
	}

	go func() {
		defer span.End()
		defer close(out)
		defer rows.Close()

		for rows.Next() {
			user := &models.User{}

			if err := rows.Scan(
				&user.Id,
				&user.LastName,
				&user.FirstName,
				&user.MiddleName,
				&user.Email,
				&user.PersonalData.Phone,
				&user.PersonalData.Snils,
				&user.PersonalData.Passport.Number,
				&user.PersonalData.Passport.Series,
				&user.CreatedAt,
				&user.UpdatedAt,
			); err != nil {
				log.Error("failed to scan query result", sl.Err(err))
				return
			}

			log.Info("scanned user", sl.UUID("id", user.Id))

			out <- user
		}
	}()

	return out, nil
}
