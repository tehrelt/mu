//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/consumption-service/internal/config"
	"github.com/tehrelt/mu/consumption-service/internal/storage/api/accountapi"
	"github.com/tehrelt/mu/consumption-service/internal/storage/api/billingapi"
	"github.com/tehrelt/mu/consumption-service/internal/storage/api/rateapi"
	"github.com/tehrelt/mu/consumption-service/internal/storage/pg/consumptionstorage"
	"github.com/tehrelt/mu/consumption-service/internal/transport/amqp"
	tgrpc "github.com/tehrelt/mu/consumption-service/internal/transport/grpc"
	"github.com/tehrelt/mu/consumption-service/internal/usecase"
	"go.opentelemetry.io/otel/trace"
)

func New(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,

		tgrpc.New,
		amqp.New,

		usecase.New,
		wire.Bind(new(usecase.AccountProvider), new(*accountapi.Api)),
		wire.Bind(new(usecase.BillingProvider), new(*billingapi.Api)),
		wire.Bind(new(usecase.ServiceProvider), new(*rateapi.Api)),
		wire.Bind(new(usecase.ConsumptionStorage), new(*consumptionstorage.Storage)),

		rateapi.New,
		billingapi.New,
		accountapi.New,

		_amqp,

		consumptionstorage.New,
		_pg,

		_tracer,
		config.New,
	))
}
func _pg(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, func(), error) {
	pool, err := pgxpool.Connect(ctx, cfg.Postgres.ConnectionString())
	if err != nil {
		return nil, nil, err
	}

	log := slog.With(slog.String("connection", cfg.Postgres.ConnectionString()))
	log.Debug("connecting to database")
	t := time.Now()
	if err := pool.Ping(ctx); err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		return nil, func() { pool.Close() }, err
	}
	log.Info("connected to database", slog.String("ping", fmt.Sprintf("%2.fs", time.Since(t).Seconds())))

	return pool, func() { pool.Close() }, nil
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
			defer conn.Close()
		}, err
	}

	closefn := func() {
		defer conn.Close()
		defer channel.Close()
	}

	exchange := cfg.ConnectServiceExchange.Exchange
	queueName := cfg.ConnectServiceExchange.Queue
	rk := cfg.ConnectServiceExchange.Routing
	log := slog.With(slog.String("exchange", exchange))
	log.Info("declaring exchange")
	if err := channel.ExchangeDeclare(exchange, "fanout", true, false, false, false, nil); err != nil {
		slog.Error("failed to declare exchange", sl.Err(err))
		return nil, closefn, err
	}

	log.Info("declaring queue", slog.String("queue", queueName))
	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Error("failed to declare queue", sl.Err(err), slog.String("queue", queueName))
		return nil, closefn, err
	}

	log.Info("binding queue", slog.String("queue", queueName))
	if err := channel.QueueBind(queue.Name, rk, exchange, false, nil); err != nil {
		log.Error("failed to bind queue", sl.Err(err), slog.String("queue", queueName))
		return nil, closefn, err
	}

	return channel, closefn, nil
}

func _servers(g *tgrpc.Server, c *amqp.AmqpConsumer) []Server {
	return []Server{g, c}
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	jaeger := cfg.Jaeger.Endpoint
	appname := cfg.App.Name

	slog.Debug("connecting to jaeger", slog.String("jaeger", jaeger), slog.String("appname", appname))

	return tracer.SetupTracer(ctx, jaeger, appname)
}
