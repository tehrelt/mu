package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/consumptionpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create implements consumptionpb.ConsumptionServiceServer.
func (s *Server) Create(ctx context.Context, req *consumptionpb.CreateRequest) (*consumptionpb.CreateResponse, error) {
	accId, err := uuid.Parse(req.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id=%s", req.AccountId)
	}

	svcId, err := uuid.Parse(req.ServiceId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id=%s", req.ServiceId)
	}

	in := &dto.NewCabinet{
		AccountId: accId,
		ServiceId: svcId,
	}

	cabinet, err := s.uc.CreateCabinet(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create cabinet: %v", err)
	}

	return &consumptionpb.CreateResponse{
		Id: cabinet.Id.String(),
	}, nil
}
