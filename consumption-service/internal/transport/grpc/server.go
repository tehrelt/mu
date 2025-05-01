package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/consumption-service/internal/config"
	"github.com/tehrelt/mu/consumption-service/internal/usecase"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/consumptionpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var _ consumptionpb.ConsumptionServiceServer = (*Server)(nil)

type Server struct {
	cfg    *config.Config
	uc     *usecase.UseCase
	tracer trace.Tracer

	consumptionpb.UnimplementedConsumptionServiceServer
}

// FindConsumption implements consumptionpb.ConsumptionServiceServer.
func (s *Server) FindConsumption(context.Context, *consumptionpb.FindConsumptionRequest) (*consumptionpb.FindConsumptionResponse, error) {
	panic("unimplemented")
}

// ListCabinets implements consumptionpb.ConsumptionServiceServer.
func (s *Server) ListCabinets(*consumptionpb.ListCabinetsRequest, grpc.ServerStreamingServer[consumptionpb.Cabinet]) error {
	panic("unimplemented")
}

func New(cfg *config.Config, uc *usecase.UseCase, t trace.Tracer) *Server {
	return &Server{
		cfg:    cfg,
		uc:     uc,
		tracer: t,
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

	consumptionpb.RegisterConsumptionServiceServer(server, s)

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
