package grpc

import (
	"log/slog"
	"math"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/consumption-service/internal/dto"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/consumptionpb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Logs implements consumptionpb.ConsumptionServiceServer.
func (s *Server) Logs(req *consumptionpb.LogsRequest, stream grpc.ServerStreamingServer[consumptionpb.LogsResponse]) error {

	span := trace.SpanFromContext(stream.Context())

	filters := &dto.LogsFilters{}

	if req.CabinetId != "" {
		uuid, err := uuid.Parse(req.CabinetId)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "bad cabinet_id=%s", req.CabinetId)
		}

		filters.CabinetId = uuid
	}

	if req.Pagination != nil {
		p := req.Pagination
		filters.Pagination.Limit = p.Limit
		filters.Pagination.Offset = p.Offset
	}

	logs, total, err := s.uc.Logs(stream.Context(), filters)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get logs: %v", err)
	}

	slog.Debug("receieved logs", slog.Any("log", logs), slog.Any("total", total))

	batchsize := 20

	span.AddEvent("sending meta response",
		trace.WithAttributes(
			attribute.Int64("total", int64(total)),
			attribute.Int("batchsize", batchsize),
		),
	)
	stream.Send(&consumptionpb.LogsResponse{
		Meta: &consumptionpb.LogsResponseMeta{
			Total:     total,
			Batchsize: uint32(batchsize),
		},
	})
	q := float64(len(logs)) / float64(batchsize)
	iters := int(math.Ceil(q))

	slog.Debug(
		"iters count",
		slog.Int("len", len(logs)),
		slog.Int("batchsize", batchsize),
		slog.Int("iters", iters),
		slog.Float64("q", q),
	)

	for i := range iters {
		start := i * batchsize
		end := min(start+batchsize, len(logs))

		batch := make([]*consumptionpb.Consumption, 0, end)

		for j := start; j < end; j++ {
			l := logs[j]

			batch = append(batch, &consumptionpb.Consumption{
				Id:        l.Id.String(),
				Consumed:  l.Consumed,
				CabinetId: l.CabinetId.String(),
				AccountId: l.AccountId.String(),
				ServiceId: l.ServiceId.String(),
				CreatedAt: l.CreatedAt.Unix(),
			})
		}

		span.AddEvent("sending batch", trace.WithAttributes(
			attribute.Int("start", start),
			attribute.Int("end", end),
		))
		stream.Send(&consumptionpb.LogsResponse{
			Consumptions: batch,
		})
	}

	return nil
}
