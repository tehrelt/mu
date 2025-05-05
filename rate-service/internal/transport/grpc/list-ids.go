package grpc

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/rate-service/pkg/pb/ratepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListIds implements ratepb.RateServiceServer.
func (s *Server) ListIds(req *ratepb.ListIdsRequest, stream grpc.ServerStreamingServer[ratepb.Service]) error {

	ids := make([]uuid.UUID, 0, len(req.Ids))
	for _, id := range req.Ids {
		parsed, err := uuid.Parse(id)
		if err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		ids = append(ids, parsed)
	}

	rates, err := s.storage.ListIds(stream.Context(), ids)
	if err != nil {
		return err
	}

	for rate := range rates {
		if err := stream.Send(&ratepb.Service{
			Id:          rate.Id,
			Name:        rate.Name,
			MeasureUnit: rate.MeasureUnit,
			Rate:        rate.Rate,
			Type:        rate.Type.ToProto(),
		}); err != nil {
			return err
		}
	}

	return nil
}
