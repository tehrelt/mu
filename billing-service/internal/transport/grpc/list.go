package grpc

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/billing-service/internal/dto"
	"github.com/tehrelt/mu/billing-service/internal/models"
	"github.com/tehrelt/mu/billing-service/pkg/pb/billingpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) List(in *billingpb.ListRequest, stream grpc.ServerStreamingServer[billingpb.ListResponse]) error {

	paymentsChannel := make(chan models.Payment, 4)
	errChannel := make(chan error)

	filters := dto.NewPaymentFilter()

	if in.GetAccountId() != "" {
		id, err := uuid.Parse(in.AccountId)
		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		filters = filters.SetAccountId(id)
	}

	if in.AmountRange != nil {
		if in.AmountRange.Min != 0 {
			filters = filters.SetAmountMin(in.AmountRange.Min)
		}
		if in.AmountRange.Max != 0 {
			filters = filters.SetAmountMax(in.AmountRange.Max)
		}
	}

	if in.Pagination != nil {
		if in.Pagination.Limit != 0 {
			filters = filters.SetLimit(in.Pagination.Limit)
		}
		if in.Pagination.Offset != 0 {
			filters = filters.SetOffset(in.Pagination.Offset)
		}
	}

	if in.Status != billingpb.PaymentStatus_nil {
		p := models.PaymentStatusNil
		st, err := p.FromProto(in.Status)
		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		filters = filters.SetStatus(st)
	}

	go func() {
		if err := s.storage.List(stream.Context(), filters, paymentsChannel); err != nil {
			errChannel <- err
		}
	}()

	go func() {
		defer close(errChannel)
		for payment := range paymentsChannel {
			res := &billingpb.ListResponse{
				Payment: payment.ToProto(),
			}

			slog.Debug("sending response", slog.Any("res", res))

			if err := stream.Send(res); err != nil {
				errChannel <- status.Error(codes.Internal, err.Error())
			}
		}
	}()

	for err := range errChannel {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
