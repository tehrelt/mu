package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/tehrelt/mu/register-service/internal/config"
	"github.com/tehrelt/mu/register-service/internal/lib/tracer/interceptors"
	"github.com/tehrelt/mu/register-service/pkg/pb/authpb"
	"github.com/tehrelt/mu/register-service/pkg/pb/registerpb"
	"github.com/tehrelt/mu/register-service/pkg/pb/userpb"
	"github.com/tehrelt/mu/register-service/pkg/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg     *config.Config
	authApi authpb.AuthServiceClient
	userApi userpb.UserServiceClient

	registerpb.UnimplementedRegisterServiceServer
}

func New(cfg *config.Config, authApi authpb.AuthServiceClient, userApi userpb.UserServiceClient) *Server {
	return &Server{
		cfg:     cfg,
		authApi: authApi,
		userApi: userApi,
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptors.UnaryServerInterceptor()),
		grpc.StreamInterceptor(interceptors.StreamServerInterceptor()),
	)
	host := s.cfg.Grpc.Host
	port := s.cfg.Grpc.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Info("start grpc server", slog.String("addr", addr))

	slog.Info("enabling reflection")
	reflection.Register(server)

	registerpb.RegisterRegisterServiceServer(server, s)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error(
			"failed to listen addr",
			slog.String("addr", addr),
			sl.Err(err),
		)
		return err
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			slog.Error(
				"failed to serve grpc server",
				sl.Err(err),
			)
		}
	}()

	<-ctx.Done()
	slog.Info("grpc server stopped")
	server.GracefulStop()
	return nil
}
