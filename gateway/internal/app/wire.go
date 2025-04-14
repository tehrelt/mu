//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/gateway/internal/config"
	"github.com/tehrelt/mu/gateway/internal/transport/http"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/registerpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/jackc/pgx/stdlib"
)

func New(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		_servers,
		http.New,

		// _userpb,
		_regpb,
		_authpb,
		_tracer,
		config.New,
	))
}

// func _userpb(cfg *config.Config) (userpb.UserServiceClient, error) {

// 	host := cfg.UserService.Host
// 	port := cfg.UserService.Port

// 	return create_grpc_client(host, port, userpb.NewUserServiceClient)
// }

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

func _regpb(cfg *config.Config) (registerpb.RegisterServiceClient, error) {
	host := cfg.RegisterService.Host
	port := cfg.RegisterService.Port

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

	return registerpb.NewRegisterServiceClient(conn), nil
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	return tracer.SetupTracer(ctx, cfg.Jaeger.Endpoint, cfg.App.Name)
}

// func create_grpc_client[T any](host string, port int, fn func(grpc.ClientConnInterface) T) (T, error) {
// 	var zero T
// 	addr := fmt.Sprintf("%s:%d", host, port)
// 	conn, err := grpc.NewClient(
// 		addr,
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()),
// 		grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
// 	)
// 	if err != nil {
// 		return zero, err
// 	}

// 	return fn(conn), nil
// }

func _servers(h *http.Server) ([]Server, error) {
	servers := []Server{h}
	return servers, nil
}
