package grpc

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu/rate-service/internal/models"
	"github.com/tehrelt/mu/rate-service/pkg/pb/ratepb"
	"github.com/tehrelt/mu/rate-service/pkg/sl"
)

// Create implements ratepb.RateServiceServer.
func (s *Server) Create(ctx context.Context, in *ratepb.CreateRequest) (*ratepb.CreateResponse, error) {

	id, err := s.storage.Create(ctx, &models.CreateService{
		Name:        in.Name,
		MeasureUnit: in.MeasureUnit,
		Rate:        in.InitialRate,
		Type:        models.ServiceTypeFromProto(in.Type),
	})
	if err != nil {
		slog.Error("failed to create service", sl.Err(err))
		return nil, err
	}

	return &ratepb.CreateResponse{
		Id: id.String(),
	}, nil

}
