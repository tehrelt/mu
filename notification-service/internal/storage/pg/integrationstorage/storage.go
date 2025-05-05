package integrationstorage

import (
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/tehrelt/mu-lib/sl"
)

type Storage struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool:   pool,
		logger: slog.With(sl.Module("integration_storage")),
	}
}
