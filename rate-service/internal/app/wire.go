//go:build wireinject
// +build wireinject

package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/rate-service/internal/config"
	"github.com/tehrelt/mu/rate-service/internal/storage/amqp/rmq"
	"github.com/tehrelt/mu/rate-service/internal/storage/pg/servicestorage"
	"github.com/tehrelt/mu/rate-service/internal/transport/grpc"
	"github.com/tehrelt/mu/rate-service/pkg/sl"

	_ "github.com/jackc/pgx/stdlib"
)

func New() (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,

		grpc.New,
		servicestorage.New,

		_rmq,
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

func _rmq(cfg *config.Config) (*rmq.RabbitMq, func(), error) {

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Amqp.User, cfg.Amqp.Pass, cfg.Amqp.Host, cfg.Amqp.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		slog.Error("failed to connect rabbitmq", sl.Err(err))
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("failed to get channel", sl.Err(err))
		return nil, func() {
			conn.Close()
		}, err
	}

	if err := ch.ExchangeDeclare(cfg.Amqp.RateChangedExchange, "direct", true, false, false, false, nil); err != nil {
		slog.Error("failed to declare exchange", sl.Err(err))
		return nil, func() {
			defer conn.Close()
			defer ch.Close()
		}, err
	}
	slog.Info("exchange declared", slog.String("exchange", cfg.Amqp.RateChangedExchange))

	q, err := ch.QueueDeclare(
		cfg.Amqp.RateChangedQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error("failed to declare queue", sl.Err(err))
		return nil, func() {
			ch.Close()
			conn.Close()
		}, err
	}
	slog.Info("queue declared", slog.String("queue", cfg.Amqp.RateChangedQueue))

	if err := ch.QueueBind(q.Name, cfg.Amqp.RateChangedRoutingKey, cfg.Amqp.RateChangedExchange, false, nil); err != nil {
		slog.Error("failed to bind queue to exchange", sl.Err(err))
		return nil, func() {
			defer conn.Close()
			defer ch.Close()
		}, err
	}
	slog.Info("queue bound to exchange", slog.String("queue", q.Name), slog.String("exchange", cfg.Amqp.RateChangedExchange))

	r := rmq.New(cfg, ch, &q)

	return r, func() {
		defer conn.Close()
		defer ch.Close()
	}, nil
}
