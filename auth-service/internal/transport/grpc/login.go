package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/services"
	"github.com/tehrelt/mu/auth-service/internal/transport/grpc/converters"
	"github.com/tehrelt/mu/auth-service/pkg/pb/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login implements authpb.AuthServiceServer.
func (s *Server) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {

	email := req.GetEmail()
	if email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email works only at this moment")
	}

	tokens, err := s.authservice.Login(ctx, &dto.LoginUser{
		Email:    email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		}

		slog.Error("failed to login", sl.Err(err))
		return nil, unexpectedError
	}

	return &authpb.LoginResponse{
		Tokens: converters.TokenPairToPb(tokens),
	}, nil
}
