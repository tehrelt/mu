package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/config"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/authpb"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var _ authpb.AuthServiceServer = (*Server)(nil)

type AuthService interface {
	Authorize(ctx context.Context, token string, roles ...models.Role) (*dto.UserClaims, error)
	Login(ctx context.Context, req *dto.LoginUser) (*dto.TokenPair, error)
	Register(ctx context.Context, req *dto.RegisterUser) (*dto.TokenPair, error)
	Refresh(ctx context.Context, userId uuid.UUID, token string) (*dto.TokenPair, error)
	Logout(ctx context.Context, userId uuid.UUID) error
}

type ProfileService interface {
	Profile(ctx context.Context, token string) (*dto.Profile, error)
}

type Server struct {
	cfg            *config.Config
	authservice    AuthService
	profileservice ProfileService
	authpb.UnimplementedAuthServiceServer
}

func New(
	cfg *config.Config,
	as AuthService,
	us ProfileService,
) *Server {
	return &Server{
		cfg:            cfg,
		authservice:    as,
		profileservice: us,
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	host := s.cfg.Grpc.Host
	port := s.cfg.Grpc.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Info("start grpc server", slog.String("addr", addr))

	log.Info("enabling reflection")
	reflection.Register(server)

	authpb.RegisterAuthServiceServer(server, s)

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
	log.Info("grpc server stopped")
	server.GracefulStop()
	return nil
}
