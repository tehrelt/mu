package userstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu/user-service/internal/models"
	"github.com/tehrelt/mu/user-service/internal/storage"
	"github.com/tehrelt/mu/user-service/internal/storage/pg"
	"github.com/tehrelt/mu/user-service/pkg/sl"
)

func (s *UserStorage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	log := slog.With(sl.Method("userstorage.userByEmail"))

	query, args, err := sq.Select("u.id, u.last_name, u.first_name, u.middle_name, u.email, pd.phone, pd.snils, pd.passport_number, pd.passport_series, u.created_at, u.updated_at").
		From(fmt.Sprintf("%s u", pg.USERS)).
		CrossJoin(fmt.Sprintf("%s pd", pg.PERSONAL_DATA)).
		Where(sq.Eq{"u.email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Warn("failed to build query", sl.Err(err))
		return nil, err
	}

	user := new(models.User)

	row := s.db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug("user not found", slog.String("email", email))
			return nil, storage.ErrUserNotFound
		}
	}

	return user, nil
}
