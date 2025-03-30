package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/account-service/internal/config"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/internal/storage/pg/accountstorage"
	"github.com/tehrelt/mu/account-service/internal/storage/rmq"
	"github.com/tehrelt/mu/account-service/pkg/pb/accountpb"
	"github.com/tehrelt/mu/account-service/pkg/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg     *config.Config
	storage *accountstorage.AccountStorage
	broker  *rmq.Broker

	accountpb.UnimplementedAccountServiceServer
}

// Create implements accountpb.AccountServiceServer.
func (s *Server) Create(ctx context.Context, in *accountpb.CreateRequest) (*accountpb.CreateResponse, error) {

	uId, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, err
	}

	hId, err := uuid.Parse(in.HouseId)
	if err != nil {
		return nil, err
	}

	accId, err := s.storage.Create(ctx, &dto.CreateAccount{
		UserId:  uId,
		HouseId: hId,
	})
	if err != nil {
		return nil, err
	}

	return &accountpb.CreateResponse{
		Id: accId.String(),
	}, nil
}

// List implements accountpb.AccountServiceServer.
func (s *Server) List(*accountpb.ListRequest, grpc.ServerStreamingServer[accountpb.Account]) error {
	panic("unimplemented")
}

// ListUsersAccounts implements accountpb.AccountServiceServer.
func (s *Server) ListUsersAccounts(*accountpb.ListUsersAccountsRequest, grpc.ServerStreamingServer[accountpb.Account]) error {
	panic("unimplemented")
}

func New(cfg *config.Config, s *accountstorage.AccountStorage, b *rmq.Broker) *Server {
	return &Server{
		cfg:     cfg,
		storage: s,
		broker:  b,
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	host := s.cfg.Grpc.Host
	port := s.cfg.Grpc.Port
	addr := fmt.Sprintf("%s:%d", host, port)
	slog.Info("start grpc server", slog.String("addr", addr))

	slog.Info("enabling reflection")
	reflection.Register(server)

	accountpb.RegisterAccountServiceServer(server, s)

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
