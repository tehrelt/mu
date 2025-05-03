//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/notification-service/internal/config"
	"github.com/tehrelt/mu/notification-service/internal/storage/pg/integrationstorage"
	"github.com/tehrelt/mu/notification-service/internal/storage/redis/otpstorage"
	"github.com/tehrelt/mu/notification-service/internal/storage/rmq"
	"github.com/tehrelt/mu/notification-service/internal/transport/amqp"
	tgrpc "github.com/tehrelt/mu/notification-service/internal/transport/grpc"
	"github.com/tehrelt/mu/notification-service/internal/usecase"
	"github.com/tehrelt/mu/notification-service/pkg/pb/ticketpb"
	"github.com/tehrelt/mu/notification-service/pkg/pb/userpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
)

func New(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,

		amqp.New,
		tgrpc.New,

		usecase.New,

		rmq.New,

		otpstorage.NewStorage,
		integrationstorage.NewStorage,

		_userpb,
		_ticketpb,
		_amqp,
		_pg,
		_redis,
		_tracer,
		config.New,
	))
}

func _redis(ctx context.Context, cfg *config.Config) (*redis.Client, func(), error) {

	log := slog.With()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Pass,
		DB:       0,
	})

	log.Debug("connecting to redis", slog.String("conn", cfg.Redis.ConnectionString()))
	t := time.Now()
	if err := client.Ping(ctx).Err(); err != nil {
		slog.Error("failed to connect to redis", slog.String("err", err.Error()), slog.String("conn", cfg.Redis.ConnectionString()))
		return nil, func() { client.Close() }, err
	}
	log.Info("connected to redis", slog.String("ping", fmt.Sprintf("%2.fs", time.Since(t).Seconds())))

	return client, func() { client.Close() }, nil
}
func _pg(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, func(), error) {

	slog.Info("connecting to database", slog.String("connection", cfg.Postgres.ConnectionString()))

	pool, err := pgxpool.Connect(ctx, "postgres://"+cfg.Postgres.ConnectionString())
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

	log := slog.With(slog.String("cfg", "_amqp"))

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

	exchange := cfg.TicketStatusChangedExchange.Exchange
	log.Info("declaring exchange", slog.String("exchange", exchange))
	if err := channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		slog.Error("failed to declare exchange", slog.String("exchange", exchange), sl.Err(err))
		return nil, closefn, err
	}

	queue, err := channel.QueueDeclare(amqp.TicketStatusChangedQueue, true, false, false, false, nil)
	if err != nil {
		return nil, closefn, err
	}

	if err := channel.QueueBind(queue.Name, cfg.TicketStatusChangedExchange.WildcardRoute, exchange, false, nil); err != nil {
		return nil, closefn, err
	}

	exchange = cfg.NotificationSendExchange.Exchange
	log.Info("declaring exchange", slog.String("exchange", exchange))
	if err := channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil); err != nil {
		slog.Error("failed to declare exchange", slog.String("exchange", exchange), sl.Err(err))
		return nil, closefn, err
	}

	slog.Info("connected to amqp", slog.String("conn", cs))

	return channel, closefn, nil
}

func _servers(g *tgrpc.Server, c *amqp.AmqpConsumer) []Server {
	return []Server{c, g}
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	jaeger := cfg.Jaeger.Endpoint
	appname := cfg.App.Name

	slog.Debug("connecting to jaeger", slog.String("jaeger", jaeger), slog.String("appname", appname))

	return tracer.SetupTracer(ctx, jaeger, appname)
}

func _ticketpb(cfg *config.Config) (
	ticketpb.TicketServiceClient,
	func(),
	error,
) {
	host := cfg.TicketService.Host
	port := cfg.TicketService.Port
	addr := fmt.Sprintf("%s:%d", host, port)

	client, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, nil, err
	}

	return ticketpb.NewTicketServiceClient(client), func() { client.Close() }, nil
}

func _userpb(cfg *config.Config) (
	userpb.UserServiceClient,
	func(),
	error,
) {
	addr := cfg.UserService.Address()

	client, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, nil, err
	}

	return userpb.NewUserServiceClient(client), func() { client.Close() }, nil
}
