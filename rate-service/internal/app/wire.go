//go:build wireinject
// +build wireinject

package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/tehrelt/mu/rate-service/internal/config"
	"github.com/tehrelt/mu/rate-service/internal/storage/pg/servicestorage"
	"github.com/tehrelt/mu/rate-service/internal/transport/grpc"

	_ "github.com/jackc/pgx/stdlib"
)

func New() (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,

		grpc.New,
		servicestorage.New,

		_pg,
		config.New,
	))
}

func _pg(cfg *config.Config) (*sqlx.DB, func(), error) {
	host := cfg.Postgres.Host
	port := cfg.Postgres.Port
	user := cfg.Postgres.User
	pass := cfg.Postgres.Pass
	name := cfg.Postgres.Name

	cs := fmt.Sprintf(`postgres://%s:%s@%s:%d/%s?sslmode=disable`, user, pass, host, port, name)

	db, err := sqlx.Connect("pgx", cs)
	if err != nil {
		return nil, nil, err
	}

	slog.Debug("connecting to database", slog.String("conn", cs))
	t := time.Now()
	if err := db.Ping(); err != nil {
		slog.Error("failed to connect to database", slog.String("err", err.Error()), slog.String("conn", cs))
		return nil, func() { db.Close() }, err
	}
	slog.Info("connected to database", slog.String("ping", fmt.Sprintf("%2.fs", time.Since(t).Seconds())))

	return db, func() { db.Close() }, nil
}

func _servers(g *grpc.Server) []Server {
	return []Server{g}
}
