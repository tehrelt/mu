package grpc

import (
	"context"

	"github.com/samber/lo"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/transport/grpc/converters"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/authpb"
)

// Profile implements authpb.AuthServiceServer.
func (s *Server) Profile(ctx context.Context, req *authpb.ProfileRequest) (*authpb.ProfileResponse, error) {

	profile, err := s.profileservice.Profile(ctx, req.AccessToken)
	if err != nil {
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
