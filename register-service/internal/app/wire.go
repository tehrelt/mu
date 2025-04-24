//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/register-service/internal/config"
	"github.com/tehrelt/mu/register-service/internal/cron"
	"github.com/tehrelt/mu/register-service/internal/services/registerservice"
	tgrpc "github.com/tehrelt/mu/register-service/internal/transport/grpc"
	"github.com/tehrelt/mu/register-service/pkg/pb/authpb"
	"github.com/tehrelt/mu/register-service/pkg/pb/userpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/jackc/pgx/stdlib"
)

func New(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,

		cron.New,
		tgrpc.New,

		registerservice.New,

		_userpb,
		_authpb,
		_tracer,
		config.New,
	))
}

func _servers(g *tgrpc.Server) []Server {
	return []Server{g}
}

func _userpb(cfg *config.Config) (userpb.UserServiceClient, error) {

	host := cfg.UserService.Host
	port := cfg.UserService.Port

	addr := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
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

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return authpb.NewAuthServiceClient(conn), nil
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	return tracer.SetupTracer(ctx, cfg.Jaeger.Endpoint, cfg.App.Name)
}
