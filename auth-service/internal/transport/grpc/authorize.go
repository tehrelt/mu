package grpc

import (
	"context"
	"log/slog"

	"github.com/samber/lo"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/transport/grpc/converters"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/authpb"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

// Authorize implements authpb.AuthServiceServer.
func (s *Server) Authorize(ctx context.Context, req *authpb.AuthorizeRequest) (*authpb.AuthorizeResponse, error) {

	log := slog.With(sl.Method("grpcserver.Authorize"))

	roles := lo.Map(req.Roles, func(r authpb.Role, _ int) models.Role {
		return converters.RoleFromPb(r)
	})

	if _, err := s.authservice.Authorize(ctx, req.Token, roles...); err != nil {
		log.Error("authorize error", sl.Err(err))
		return nil, err
	}

	return &authpb.AuthorizeResponse{}, nil
}
