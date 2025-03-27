package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/billing-service/internal/dto"
	"github.com/tehrelt/mu/billing-service/internal/models"
	"github.com/tehrelt/mu/billing-service/internal/storage"
	"github.com/tehrelt/mu/billing-service/pkg/pb/billingpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Cancel(ctx context.Context, in *billingpb.CancelRequest) (*billingpb.CancelResponse, error) {

	id, err := uuid.Parse(in.PaymentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	newdata := dto.NewUpdatePayment(id).SetStatus(models.PaymentStatusCanceled)

	old, err := s.storage.Update(ctx, newdata)
	if err != nil {
		if errors.Is(err, storage.ErrPaymentNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := s.broker.PublishStatusChanged(ctx, &dto.EventStatusChanged{
		PaymentId: id.String(),
		AccountId: old.AccountId,
		NewStatus: models.PaymentStatusCanceled,
	}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &billingpb.CancelResponse{}, nil
}
