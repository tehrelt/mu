package grpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/rate-service/internal/models"
	"github.com/tehrelt/mu/rate-service/internal/storage"
	"github.com/tehrelt/mu/rate-service/pkg/pb/ratepb"
	"github.com/tehrelt/mu/rate-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateServiceRate implements ratepb.RateServiceServer.
func (s *Server) UpdateServiceRate(ctx context.Context, in *ratepb.UpdateServiceRateRequest) (*ratepb.UpdateServiceRateResponse, error) {

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	service, err := s.storage.Find(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrServiceNotFound) {
			slog.Warn("no service for update found", sl.UUID("service_id", id))
			return nil, status.Errorf(codes.NotFound, "no service to update found")
		}
		slog.Error("failed to find service", sl.Err(err), sl.UUID("service_id", id))
		return nil, err
	}

	if err := s.storage.Update(ctx, &models.UpdateServiceRate{
		Id:   id,
		Rate: in.NewRate,
	}); err != nil {
		if errors.Is(err, storage.ErrServiceNotFound) {
			slog.Warn("no service for update found", sl.UUID("service_id", id))
			return nil, status.Errorf(codes.NotFound, "no service to update found")
		}
		slog.Error("failed to update rate of service", sl.Err(err), sl.UUID("service_id", id))
		return nil, err
	}

	if err := s.broker.NotifyRateChanged(ctx, &models.EventRateChanged{
		Id:        id,
		NewRate:   in.NewRate,
		OldRate:   service.Rate,
		Timestamp: time.Now(),
	}); err != nil {
		slog.Error("failed to notify about rate change", sl.Err(err), sl.UUID("service_id", id))
		return nil, err
	}

	return &ratepb.UpdateServiceRateResponse{}, nil
}
