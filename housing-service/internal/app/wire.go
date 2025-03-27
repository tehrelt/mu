//go:build wireinject
// +build wireinject

package app

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu/housing-service/internal/config"
	"github.com/tehrelt/mu/housing-service/internal/storage/pg/housestorage"
	"github.com/tehrelt/mu/housing-service/internal/storage/rmq"
	"github.com/tehrelt/mu/housing-service/internal/transport/amqp"
	tgrpc "github.com/tehrelt/mu/housing-service/internal/transport/grpc"
	ratepb "github.com/tehrelt/mu/housing-service/pkg/pb/ratespb"
	"github.com/tehrelt/mu/housing-service/pkg/sl"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/jackc/pgx/stdlib"
)

func New() (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,

		tgrpc.New,
		amqp.New,

		housestorage.New,
		rmq.New,

		_ratepb,
		_amqp,
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

func _amqp(cfg *config.Config) (*amqp091.Channel, func(), error) {
	cs := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Amqp.User, cfg.Amqp.Pass, cfg.Amqp.Host, cfg.Amqp.Port)

	conn, err := amqp091.Dial(cs)
	if err != nil {
		return nil, nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, func() {
			conn.Close()
		}, err
	}

	closefn := func() {
		defer conn.Close()
		defer channel.Close()
	}

	if err := amqp_setup_exchange(channel, cfg.ConnectServiceQueue.Exchange, cfg.ConnectServiceQueue.Routing); err != nil {
		slog.Error("failed to setup notifications exchange", sl.Err(err))
		return nil, closefn, err
	}

	if err := amqp_setup_exchange(channel, cfg.ServiceConnectedQueue.Exchange, cfg.ServiceConnectedQueue.Routing); err != nil {
		slog.Error("failed to setup notifications exchange", sl.Err(err))
		return nil, closefn, err
	}

	return channel, closefn, nil
}
func amqp_setup_exchange(channel *amqp091.Channel, exchange string, queues ...string) error {

	log := slog.With(slog.String("exchange", exchange))
	log.Info("declaring exchange")
	if err := channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		slog.Error("failed to declare notifications queue", sl.Err(err))
		return err
	}

	for _, queueName := range queues {
		log.Info("declaring queue", slog.String("queue", queueName))
		queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
		if err != nil {
			log.Error("failed to declare queue", sl.Err(err), slog.String("queue", queueName))
			return err
		}

		log.Info("binding queue", slog.String("queue", queueName))
		if err := channel.QueueBind(queue.Name, queueName, exchange, false, nil); err != nil {
			log.Error("failed to bind queue", sl.Err(err), slog.String("queue", queueName))
			return err
		}
	}

	return nil
}

func _servers(g *tgrpc.Server, c *amqp.AmqpConsumer) []Server {
	return []Server{g, c}
}

func _ratepb(cfg *config.Config) (ratepb.RateServiceClient, func(), error) {
	host := cfg.RateService.Host
	port := cfg.RateService.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Debug("connecting to rate service", slog.String("addr", addr))
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := ratepb.NewRateServiceClient(conn)
	return client, func() { conn.Close() }, nil
}
