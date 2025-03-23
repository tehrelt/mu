package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/user-service/internal/config"
	"github.com/tehrelt/moi-uslugi/user-service/internal/models"
	"github.com/tehrelt/moi-uslugi/user-service/pkg/pb/userspb"
	"github.com/tehrelt/moi-uslugi/user-service/pkg/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var _ userspb.UserServiceServer = (*Server)(nil)

type UserProvider interface {
	UserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserCreator interface {
	Create(ctx context.Context, req *models.CreateUser) (uuid.UUID, error)
}

type Users struct {
	creator  UserCreator
	provider UserProvider
}

type Server struct {
	cfg   *config.Config
	users Users

	userspb.UnimplementedUserServiceServer
}

func New(cfg *config.Config, usersCreator UserCreator, usersProvider UserProvider) *Server {
	return &Server{
		cfg: cfg,
		users: Users{
			creator:  usersCreator,
			provider: usersProvider,
		},
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

	userspb.RegisterUserServiceServer(server, s)

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
