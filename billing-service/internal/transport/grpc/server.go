package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/billing-service/internal/config"
	"github.com/tehrelt/mu/billing-service/internal/storage/pg/paymentstorage"
	"github.com/tehrelt/mu/billing-service/internal/storage/rmq"
	"github.com/tehrelt/mu/billing-service/pkg/pb/billingpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg     *config.Config
	storage *paymentstorage.PaymentStorage
	broker  *rmq.Broker
	tracer  trace.Tracer

	billingpb.UnimplementedBillingServiceServer
}

func New(cfg *config.Config, s *paymentstorage.PaymentStorage, b *rmq.Broker, t trace.Tracer) *Server {
	return &Server{
		cfg:     cfg,
		storage: s,
		broker:  b,
		tracer:  t,
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

	billingpb.RegisterBillingServiceServer(server, s)

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
