package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/account-service/pkg/pb/accountpb"
	"github.com/tehrelt/mu/account-service/pkg/pb/housepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Find(ctx context.Context, in *accountpb.FindRequest) (*accountpb.FindResponse, error) {

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	account, err := s.storage.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	house, err := s.houser.Find(ctx, &housepb.HouseRequest{
		HouseId: account.HouseId,
	})
	if err != nil {
		return nil, err
	}

	return &accountpb.FindResponse{
		Account: account.ToProto(house.House),
	}, nil
}
