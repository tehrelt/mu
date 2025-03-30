package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/tehrelt/mu/register-service/internal/config"
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

func (s *Server) Register(ctx context.Context, in *registerpb.RegisterRequest) (*registerpb.RegisterResponse, error) {

	user, err := s.userApi.Create(ctx, &userpb.CreateRequest{
		Fio: &userpb.FIO{
			Lastname:   in.User.LastName,
			Firstname:  in.User.FirstName,
			Middlename: in.User.MiddleName,
		},
		Email: in.User.Email,
		PersonalData: &userpb.PersonalData{
			Passport: &userpb.Passport{
				Series: in.User.PassportSeries,
				Number: in.User.PassportNumber,
			},
			Snils: in.User.Snils,
			Phone: in.User.Phone,
		},
	})
	if err != nil {
		slog.Error("failed to create user", sl.Err(err))
		return nil, err
	}

	tokens, err := s.authApi.Register(ctx, &authpb.RegisterRequest{
		UserId:   user.Id,
		Password: in.User.Password,
	})
	if err != nil {
		slog.Error("failed to register user", sl.Err(err))
		return nil, err
	}

	return &registerpb.RegisterResponse{
		Tokens: &registerpb.Tokens{
			AccessToken:  tokens.Tokens.AccessToken,
			RefreshToken: tokens.Tokens.RefreshToken,
		},
	}, nil

}

func New(cfg *config.Config, authApi authpb.AuthServiceClient, userApi userpb.UserServiceClient) *Server {
	return &Server{
		cfg:     cfg,
		authApi: authApi,
		userApi: userApi,
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
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
