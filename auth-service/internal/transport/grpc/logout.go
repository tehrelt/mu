package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/services"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/authpb"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Logout implements authpb.AuthServiceServer.
func (s *Server) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {

	log := slog.With(sl.Method("grpcserver.Logout"))

	user, err := s.authservice.Authorize(ctx, req.AccessToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			return nil, status.Errorf(codes.Unauthenticated, "")
		}

		log.Error("failed to authorize", sl.Err(err))
		return nil, unexpectedError
	}

	userId, err := uuid.Parse(user.Id)
	if err != nil {
		log.Error("failed to parse uuid", sl.Err(err))
		return nil, unexpectedError
	}

	if err := s.authservice.Logout(ctx, userId); err != nil {
		log.Error("failed to logout", sl.Err(err))
		return nil, unexpectedError
	}

	return &authpb.LogoutResponse{}, nil
}
