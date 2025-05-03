//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/notification-service/internal/config"
	"github.com/tehrelt/mu/notification-service/internal/storage/rmq"
	"github.com/tehrelt/mu/notification-service/internal/transport/amqp"
	"github.com/tehrelt/mu/notification-service/pkg/pb/ticketpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/jackc/pgx/stdlib"
)

func New(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,

		amqp.New,

		rmq.New,

		// _ticketpb,
		_amqp,
		// _pg,
		_tracer,
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

func _servers(c *amqp.AmqpConsumer) []Server {
	return []Server{c}
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
