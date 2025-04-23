package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/pkg/pb/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) AddRole(ctx context.Context, req *authpb.AddRoleRequest) (*authpb.AddRoleResponse, error) {

	userid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	var role models.Role
	role = role.FromProto(req.Role)
	if !role.Valid() {
		return nil, status.Error(codes.InvalidArgument, "invalid role")
	}

	data := &dto.UserRoles{
		UserId: userid,
		Roles:  []models.Role{role},
	}

	if err := s.roleservice.Add(ctx, data); err != nil {
		return nil, err
	}

	return &authpb.AddRoleResponse{}, nil
}
