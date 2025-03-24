package credentialstorage

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/storage/pg"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

func (s *CredentialStorage) Save(ctx context.Context, creds *models.Credentials) (err error) {
	log := slog.With(sl.Method("credentialstorage.Save"))

	log.Debug("saving credentials")
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))
		return err
	}

	defer func() {
		if err != nil {
			log.Debug("rollback tx")
			_ = tx.Rollback()
			return
		}

		log.Debug("commit tx")
		err = tx.Commit()
		if err != nil {
			log.Warn("failed to commit transaction", sl.Err(err))
		}
	}()

	query, args, err := sq.Insert(pg.CREDENTIALS_TABLE).
		Columns("id", "hashed_password").
		Values(creds.UserId, creds.HashedPassword).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Error("failed to generate query", sl.Err(err))
		return err
	}

	log.Debug("executing insert query", slog.String("query", query), slog.Any("args", args))
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Warn("failed to execute query", sl.Err(err))
		return err
	}

	if len(creds.Roles) != 0 {
		builder := sq.Insert(pg.ROLES_TABLE).Columns("user_id", "role").PlaceholderFormat(sq.Dollar)

		for _, role := range creds.Roles {
			if !role.Valid() {
				log.Warn("invalid role in request", slog.String("role", string(role)))
				continue
			}

			builder = builder.Values(creds.UserId, role)
		}

		query, args, err = builder.ToSql()
		if err != nil {
			log.Warn("failed to build query for saving roles")
			return err
		}

		log.Debug("saving roles", slog.String("query", query), slog.Any("args", args))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			log.Error("failed to execute query", sl.Err(err))
			return err
		}
	}

	return nil
}
