package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/consumptionpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Consume(ctx context.Context, req *consumptionpb.ConsumeRequest) (*consumptionpb.ConsumeResponse, error) {
	cabinetId, err := uuid.Parse(req.CabinetId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid cabinet id=%s", req.CabinetId)
	}

	created, err := s.uc.Consume(ctx, &dto.NewConsume{
		FindCabinet: dto.FindCabinet{
			Id: cabinetId,
		},
		Consumed: req.Consumed,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create consumption: %v", err)
	}

	return &consumptionpb.ConsumeResponse{
		Id: created.LogId.String(),
	}, nil
}
