package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/samber/lo"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/internal/services"
	"github.com/tehrelt/mu/auth-service/internal/transport/grpc/converters"
	"github.com/tehrelt/mu/auth-service/pkg/pb/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Profile implements authpb.AuthServiceServer.
func (s *Server) Profile(ctx context.Context, req *authpb.ProfileRequest) (*authpb.ProfileResponse, error) {

	profile, err := s.profileservice.Profile(ctx, req.AccessToken)
	if err != nil {
		slog.Error("error with get profile", sl.Err(err))
		if errors.Is(err, services.ErrTokenExpired) {
			return nil, status.Errorf(codes.Unauthenticated, "token expired")
		}
		return nil, err
	}

	roles := lo.Map(profile.Roles, func(r models.Role, _ int) authpb.Role {
		return converters.RoleToPb(r)
	})

	return &authpb.ProfileResponse{
		Id:         profile.Id.String(),
		LastName:   profile.LastName,
		FirstName:  profile.FirstName,
		MiddleName: profile.MiddleName,
		Email:      profile.Email,
		Roles:      roles,
	}, nil
}
