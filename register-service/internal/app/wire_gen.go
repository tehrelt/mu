// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"context"
	"fmt"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/register-service/internal/config"
	"github.com/tehrelt/mu/register-service/internal/cron"
	"github.com/tehrelt/mu/register-service/internal/services/registerservice"
	"github.com/tehrelt/mu/register-service/internal/transport/grpc"
	"github.com/tehrelt/mu/register-service/pkg/pb/authpb"
	"github.com/tehrelt/mu/register-service/pkg/pb/userpb"
	"go.opentelemetry.io/otel/trace"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

import (
	_ "github.com/jackc/pgx/stdlib"
)

// Injectors from wire.go:

func New(ctx context.Context) (*App, func(), error) {
	configConfig := config.New()
	userServiceClient, err := _userpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	authServiceClient, err := _authpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	service := registerservice.New(userServiceClient, authServiceClient)
	server := grpc.New(configConfig, service)
	v := _servers(server)
	tracer, err := _tracer(ctx, configConfig)
	if err != nil {
		return nil, nil, err
	}
	cronCron := cron.New(configConfig, service)
	app := newApp(v, tracer, cronCron)
	return app, func() {
	}, nil
}

// wire.go:

func _servers(g *grpc.Server) []Server {
	return []Server{g}
}

func _userpb(cfg *config.Config) (userpb.UserServiceClient, error) {

	host := cfg.UserService.Host
	port := cfg.UserService.Port

	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc2.NewClient(
		addr, grpc2.WithTransportCredentials(insecure.NewCredentials()), grpc2.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()), grpc2.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return userpb.NewUserServiceClient(conn), nil
}

func _authpb(cfg *config.Config) (authpb.AuthServiceClient, error) {

	host := cfg.AuthService.Host
	port := cfg.AuthService.Port

	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc2.NewClient(
		addr, grpc2.WithTransportCredentials(insecure.NewCredentials()), grpc2.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()), grpc2.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return authpb.NewAuthServiceClient(conn), nil
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	return tracer.SetupTracer(ctx, cfg.Jaeger.Endpoint, cfg.App.Name)
}
