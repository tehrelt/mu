package grpc

import (
	"context"

	"github.com/tehrelt/mu/billing-service/internal/dto"
	"github.com/tehrelt/mu/billing-service/pkg/pb/billingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Create(ctx context.Context, in *billingpb.CreateRequest) (*billingpb.CreateResponse, error) {

	p := &dto.CreatePayment{}
	if err := p.FromProto(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid body")
	}

	id, err := s.storage.Create(ctx, p)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &billingpb.CreateResponse{
		Id: id.String(),
	}, nil
}
