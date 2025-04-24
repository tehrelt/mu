package grpc

import (
	"context"

	"github.com/tehrelt/mu/auth-service/pkg/pb/authpb"
)

func (s *Server) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {

	tokens, err := s.authservice.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &authpb.RefreshResponse{
		Tokens: &authpb.Tokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		},
	}, nil
}
