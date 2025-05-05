//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/ticket-service/internal/config"
	"github.com/tehrelt/mu/ticket-service/internal/storage/mongo/ticketstorage"
	"github.com/tehrelt/mu/ticket-service/internal/storage/rmq"
	tgrpc "github.com/tehrelt/mu/ticket-service/internal/transport/grpc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.opentelemetry.io/otel/trace"
)

func New(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,

		tgrpc.New,

		_ticketstorage,
		rmq.New,

		_amqp,
		_mongo,
		_tracer,
		config.New,
	))
}

func _ticketstorage(cfg *config.Config, client *mongo.Client) (*ticketstorage.Storage, error) {
	db := client.Database(cfg.Mongo.Database)
	if db == nil {
		return nil, errors.New("failed to select database")
	}

	return ticketstorage.New(db), nil
}

func _mongo(ctx context.Context, cfg *config.Config) (*mongo.Client, func(), error) {
	slog.Debug("connecting to mongodb", slog.String("cs", cfg.Mongo.ConnectionString()))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.ConnectionString()))
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		defer client.Disconnect(ctx)
	}

	start := time.Now()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, cleanup, err
	}
	slog.Info("pinged mongodb", slog.Duration("ping", time.Since(start)))

	return client, cleanup, nil
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

	exchange := cfg.TicketStatusChangedExchange.Exchange
	log := slog.With(slog.String("exchange", exchange))
	log.Info("declaring exchange")
	if err := channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		slog.Error("failed to declare notifications queue", sl.Err(err))
		return nil, closefn, err
	}

	return channel, closefn, nil
}

func _servers(g *tgrpc.Server) []Server {
	return []Server{g}
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	jaeger := cfg.Jaeger.Endpoint
	appname := cfg.App.Name

	slog.Debug("connecting to jaeger", slog.String("jaeger", jaeger), slog.String("appname", appname))

	return tracer.SetupTracer(ctx, jaeger, appname)
}
