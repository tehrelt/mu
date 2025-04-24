package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/internal/services"
	"github.com/tehrelt/mu/auth-service/pkg/pb/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {

	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	roles := make([]models.Role, len(req.Roles))
	for i, role := range req.Roles {
		roles[i] = roles[i].FromProto(role)
	}

	reguser := &dto.RegisterUser{
		UserId:   uid,
		Password: req.Password,
		Roles:    roles,
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
