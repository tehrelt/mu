package grpc

import (
	"context"

	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/pkg/pb/authpb"
)

func (s *Server) HasAdmin(ctx context.Context, req *authpb.HasAdminRequest) (*authpb.HasAdminResponse, error) {

	count, err := s.roleservice.Count(ctx, models.Role_Admin)
	if err != nil {
		return nil, err
	}

	has := count > 0

	return &authpb.HasAdminResponse{
		HasAdmin: has,
	}, nil
}
