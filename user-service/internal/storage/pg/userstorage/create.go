package userstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/user-service/internal/models"
	"github.com/tehrelt/moi-uslugi/user-service/internal/storage/pg"
	"github.com/tehrelt/moi-uslugi/user-service/pkg/sl"
)

func (s *UserStorage) Create(ctx context.Context, user *models.CreateUser) (id uuid.UUID, err error) {

	log := slog.With(sl.Method("userstorage.Create"))

	log.Debug("creating user", slog.Any("create user dto", user))

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

	sql, args, err := sq.Insert(pg.USERS).
		Columns("last_name", "first_name", "middle_name", "email").
		Values(user.LastName, user.FirstName, user.MiddleName, user.Email).
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

	sql, args, err = sq.Insert(pg.PERSONAL_DATA).
		Columns("user_id", "phone", "passport_number", "passport_series", "snils").
		Values(id.String(), user.PersonalData.Phone, user.PersonalData.Passport.Number, user.PersonalData.Passport.Series, user.PersonalData.Snils).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return
	}

	log.Debug("executing query", slog.String("sql", sql), slog.Any("args", args))

	if _, err = tx.ExecContext(ctx, sql, args...); err != nil {
		log.Error("failed to execute query", sl.Err(err))
	}

	return
}
