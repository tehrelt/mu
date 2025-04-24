package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/billing-service/internal/storage"
	"github.com/tehrelt/mu/billing-service/pkg/pb/billingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Find(ctx context.Context, in *billingpb.FindRequest) (*billingpb.FindResponse, error) {

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	p, err := s.storage.Find(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrPaymentNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		return nil, err
	}

	return &billingpb.FindResponse{
		Payment: p.ToProto(),
	}, nil

}
