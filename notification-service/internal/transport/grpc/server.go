package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/notification-service/internal/config"
	"github.com/tehrelt/mu/notification-service/internal/usecase"
	"github.com/tehrelt/mu/notification-service/pkg/pb/notificationpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	cfg    *config.Config
	uc     *usecase.UseCase
	tracer trace.Tracer

	notificationpb.UnimplementedNotificationServiceServer
}

// Integrations implements notificationpb.NotificationServiceServer.
func (s *Server) Integrations(ctx context.Context, req *notificationpb.IntegrationsRequest) (*notificationpb.IntegrationsResponse, error) {

	userid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id=%s", req.UserId)
	}

	settings, err := s.uc.Find(ctx, userid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user settings")
	}

	if settings == nil {
		return nil, status.Errorf(codes.NotFound, "user settings not found")
	}

	resp := &notificationpb.IntegrationsResponse{}

	if settings.TelegramChatId != nil {
		resp.TelegramChatId = *settings.TelegramChatId
	}

	return resp, nil
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

	notificationpb.RegisterNotificationServiceServer(server, s)

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
