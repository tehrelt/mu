package grpc

import (
	"context"

	"github.com/tehrelt/mu/consumption-service/pkg/pb/consumptionpb"
)

func (s *Server) FindCabinet(ctx context.Context, req *consumptionpb.FindCabinetRequest) (*consumptionpb.FindCabinetResponse, error) {
	panic("unimplemented")
}
