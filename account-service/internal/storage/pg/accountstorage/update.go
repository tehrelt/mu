package accountstorage

import (
	"context"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/internal/models"
	"github.com/tehrelt/mu/account-service/internal/storage/pg"
	"go.opentelemetry.io/otel"
)

func (s *AccountStorage) Update(ctx context.Context, in *dto.UpdateAccount) (*models.Account, error) {

	fn := "accountstroage.Update"
	t := otel.Tracer(tracer.TracerKey)
	ctx, span := t.Start(ctx, fn)
	defer span.End()
	log := slog.With(sl.Method("accountstorage.Update"))

	log.Debug("updating account", slog.Any("update account dto", in))

	oldpayment, err := s.Find(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	builder := sq.Update(pg.ACCOUNTS_TABLE).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": in.Id}).
		Set("updated_at", time.Now())

	if in.NewBalance != nil {
		builder = builder.Set("balance", *in.NewBalance)
	}

	query, args, err := builder.
		ToSql()

	if err != nil {
		log.Error("failed to build query", sl.Err(err))
		return nil, err
	}

	log.Debug("executing query", sl.Query(query), sl.Args(args))

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		log.Error("failed to execute query", sl.Err(err))
		return nil, err
	}

	return oldpayment, nil
}
