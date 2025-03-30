package grpc

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu/register-service/pkg/pb/authpb"
	"github.com/tehrelt/mu/register-service/pkg/pb/registerpb"
	"github.com/tehrelt/mu/register-service/pkg/pb/userpb"
	"github.com/tehrelt/mu/register-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				return nil, status.Error(codes.AlreadyExists, e.Message())
			}
		}
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
