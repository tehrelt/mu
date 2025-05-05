package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/consumptionpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) FindCabinet(ctx context.Context, req *consumptionpb.FindCabinetRequest) (*consumptionpb.FindCabinetResponse, error) {

	criteria := &dto.FindCabinet{}

	switch w := req.Criteria.(type) {
	case *consumptionpb.FindCabinetRequest_Id:
		id, err := uuid.Parse(w.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid cabinet id = %s", w.Id)
		}
		criteria.Id = id
	case *consumptionpb.FindCabinetRequest_ViaAccount:
		accid, err := uuid.Parse(w.ViaAccount.AccountId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid account id = %s", w.ViaAccount)
		}

		svcid, err := uuid.Parse(w.ViaAccount.ServiceId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid service id = %s", w.ViaAccount)
		}

		criteria.AccountId = accid
		criteria.ServiceId = svcid
	}

	cabinet, err := s.uc.FindCabinet(ctx, criteria)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find cabinet: %v", err)
	}

	res := &consumptionpb.FindCabinetResponse{
		Cabinet: &consumptionpb.Cabinet{
			Id:        cabinet.Id.String(),
			AccountId: cabinet.AccountId.String(),
			ServiceId: cabinet.ServiceId.String(),
			Consumed:  cabinet.Consumed,
			CreatedAt: cabinet.CreatedAt.Unix(),
			UpdatedAt: 0,
		},
	}

	if cabinet.UpdatedAt != nil {
		res.Cabinet.UpdatedAt = cabinet.UpdatedAt.Unix()
	}

	return res, nil
}
