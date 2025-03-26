package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/services"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/authpb"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {

	reguser := &dto.RegisterUser{
		Fio: dto.Fio{
			FirstName:  req.User.FirstName,
			LastName:   req.User.LastName,
			MiddleName: req.User.MiddleName,
		},
		PersonalData: dto.PersonalData{
			Phone: req.User.Phone,
			Passport: dto.Passport{
				Number: int(req.User.PassportNumber),
				Series: int(req.User.PassportSeries),
			},
			Snils: req.User.Snils,
		},
		Email:    req.User.Email,
		Password: req.User.Password,
	}

	tokens, err := s.authservice.Register(ctx, reguser)
	if err != nil {
		if errors.Is(err, services.ErrUserExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}

		slog.Error("failed to register", sl.Err(err))
		return nil, unexpectedError
	}

	return &authpb.RegisterResponse{
		Tokens: &authpb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}
