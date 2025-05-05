//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/mailer/internal/config"
	"github.com/tehrelt/mu/mailer/internal/storage/mailer"
	"github.com/tehrelt/mu/mailer/internal/transport/amqp"
	"github.com/tehrelt/mu/mailer/internal/usecase"
	"github.com/tehrelt/mu/mailer/pkg/pb/notificationpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewApp(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		amqp.New,

		usecase.New,
		mailer.New,

		_notificationpb,

		_tracer,
		config.New,
	))
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	return tracer.SetupTracer(ctx, cfg.Jaeger.Endpoint, cfg.App.Name)
}

func _notificationpb(cfg *config.Config) (
	notificationpb.NotificationServiceClient,
	func(),
	error,
) {
	host := cfg.NotificationService.Host
	port := cfg.NotificationService.Port
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

	return notificationpb.NewNotificationServiceClient(client), func() { client.Close() }, nil
}
