package grpc

import (
	"context"

	"github.com/tehrelt/mu/register-service/pkg/pb/registerpb"
)

func (s *Server) Register(ctx context.Context, in *registerpb.RegisterRequest) (*registerpb.RegisterResponse, error) {
	return s.regservice.Register(ctx, in)
}
