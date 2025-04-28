package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/rate-service/internal/config"
	"github.com/tehrelt/mu/rate-service/internal/models"
	"github.com/tehrelt/mu/rate-service/internal/storage"
	"github.com/tehrelt/mu/rate-service/internal/storage/amqp/rmq"
	"github.com/tehrelt/mu/rate-service/internal/storage/pg/servicestorage"
	"github.com/tehrelt/mu/rate-service/pkg/pb/ratepb"
	"github.com/tehrelt/mu/rate-service/pkg/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	cfg     *config.Config
	storage *servicestorage.ServiceStorage
	broker  *rmq.RabbitMq

	ratepb.UnimplementedRateServiceServer
}

// Find implements ratepb.RateServiceServer.
func (s *Server) Find(ctx context.Context, in *ratepb.FindRequest) (*ratepb.Service, error) {

	id, err := uuid.Parse(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	service, err := s.storage.Find(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrServiceNotFound) {
			return nil, status.Errorf(codes.NotFound, "service not found")
		}
		return nil, err
	}

	return &ratepb.Service{
		Id:          service.Id,
		Name:        service.Name,
		MeasureUnit: service.MeasureUnit,
		Rate:        service.Rate,
	}, nil

}

// List implements ratepb.RateServiceServer.
func (s *Server) List(in *ratepb.ListRequest, stream grpc.ServerStreamingServer[ratepb.Service]) error {

	filters := models.NewRateFilters()

	if in.Type != ratepb.ServiceType_UNKNOWN {
		filters = filters.WithType(models.ServiceTypeFromProto(in.Type))
	}

	rates, err := s.storage.List(stream.Context(), filters)
	if err != nil {
		slog.Error("failed to list services", sl.Err(err))
		return err
	}

	for service := range rates {
		stream.Send(&ratepb.Service{
			Id:          service.Id,
			Name:        service.Name,
			MeasureUnit: service.MeasureUnit,
			Rate:        service.Rate,
			Type:        service.Type.ToProto(),
		})
	}

	return nil
}

func New(cfg *config.Config, s *servicestorage.ServiceStorage, b *rmq.RabbitMq) *Server {
	return &Server{
		cfg:     cfg,
		storage: s,
		broker:  b,
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptors.UnaryServerInterceptor()),
		grpc.StreamInterceptor(interceptors.StreamServerInterceptor()),
	)

	host := s.cfg.Grpc.Host
	port := s.cfg.Grpc.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Info("start grpc server", slog.String("addr", addr))

	slog.Info("enabling reflection")
	reflection.Register(server)

	ratepb.RegisterRateServiceServer(server, s)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error(
			"failed to listen addr",
			slog.String("addr", addr),
			sl.Err(err),
		)
		return err
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			slog.Error(
				"failed to serve grpc server",
				sl.Err(err),
			)
		}
	}()

	<-ctx.Done()
	slog.Info("grpc server stopped")
	server.GracefulStop()
	return nil
}
