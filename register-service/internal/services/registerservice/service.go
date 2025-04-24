package registerservice

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/register-service/pkg/pb/authpb"
	"github.com/tehrelt/mu/register-service/pkg/pb/registerpb"
	"github.com/tehrelt/mu/register-service/pkg/pb/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	userApi userpb.UserServiceClient
	authApi authpb.AuthServiceClient
}

func New(userApi userpb.UserServiceClient, authApi authpb.AuthServiceClient) *Service {
	return &Service{
		userApi: userApi,
		authApi: authApi,
	}
}

func (s *Service) HasAdmin(ctx context.Context) (bool, error) {
	resp, err := s.authApi.HasAdmin(ctx, &authpb.HasAdminRequest{})
	if err != nil {
		return false, err
	}

	return resp.HasAdmin, nil
}

func (s *Service) Register(ctx context.Context, in *registerpb.RegisterRequest) (*registerpb.RegisterResponse, error) {
	return s.register(ctx, in)
}

func (s *Service) register(ctx context.Context, in *registerpb.RegisterRequest) (*registerpb.RegisterResponse, error) {

	fn := "registerservice.Register"
	log := slog.With("fn", fn)

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
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return nil, status.Error(codes.AlreadyExists, e.Message())
			}
		}
		slog.Error("failed to create user", sl.Err(err))
		return nil, err
	}
	log = log.With(slog.String("userId", user.Id))
	log.Info("user created")

	roles := make([]authpb.Role, len(in.User.Roles))
	for i, role := range in.User.Roles {
		roles[i] = authpb.Role(role)
	}

	tokens, err := s.authApi.Register(ctx, &authpb.RegisterRequest{
		UserId:   user.Id,
		Password: in.User.Password,
		Roles:    roles,
	})
	if err != nil {
		slog.Error("failed to register user", sl.Err(err))
		return nil, err
	}
	log.Info("user registered")

	return &registerpb.RegisterResponse{
		Tokens: &registerpb.Tokens{
			AccessToken:  tokens.Tokens.AccessToken,
			RefreshToken: tokens.Tokens.RefreshToken,
		},
	}, nil
}
