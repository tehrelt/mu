package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/account-service/internal/config"
	"github.com/tehrelt/mu/account-service/internal/dto"
	"github.com/tehrelt/mu/account-service/internal/storage/pg/accountstorage"
	"github.com/tehrelt/mu/account-service/internal/storage/rmq"
	"github.com/tehrelt/mu/account-service/pkg/pb/accountpb"
	"github.com/tehrelt/mu/account-service/pkg/pb/housepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg     *config.Config
	storage *accountstorage.AccountStorage
	broker  *rmq.Broker

	houser housepb.HouseServiceClient

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
func (s *Server) List(in *accountpb.ListRequest, stream grpc.ServerStreamingServer[accountpb.Account]) error {
	errChan := make(chan error)
	filters := dto.NewAccountFilter()

	accChan, err := s.storage.List(stream.Context(), filters)
	if err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()

		case err := <-errChan:
			return err

		case acc, ok := <-accChan:
			if !ok {
				slog.Debug("failed to read from accounts channel")
				return nil
			}

			data := acc.ToProto()
			if err := stream.Send(data); err != nil {
				return err
			}
		}
	}
}

// ListUsersAccounts implements accountpb.AccountServiceServer.
func (s *Server) ListUsersAccounts(in *accountpb.ListUsersAccountsRequest, stream grpc.ServerStreamingServer[accountpb.Account]) error {

	fn := "grpc.ListUsersAccounts"
	log := slog.With(slog.String("fn", fn))

	filters := dto.NewAccountFilter().SetUserId(in.UserId)

	accChan, err := s.storage.List(stream.Context(), filters)
	if err != nil {
		log.Error("failed to get list of accounts", sl.Err(err))
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()

		case acc, ok := <-accChan:
			if !ok {
				log.Debug("failed to read from accounts channel")
				return nil
			}

			reqlog := log.With(slog.String("houseId", acc.HouseId), slog.String("accId", acc.Id))
			reqlog.Debug("fetching house for account")
			house, err := s.houser.House(stream.Context(), &housepb.HouseRequest{
				HouseId: acc.HouseId,
			})
			if err != nil {
				reqlog.Error("failed to fetch", sl.Err(err))
				return err
			}

			data := acc.ToProto()

			data.House = &accountpb.House{
				Id:      house.House.Id,
				Address: house.House.Address,
			}

			if err := stream.Send(data); err != nil {
				return err
			}
		}
	}
}

func New(cfg *config.Config, s *accountstorage.AccountStorage, b *rmq.Broker, h housepb.HouseServiceClient) *Server {
	return &Server{
		cfg:     cfg,
		storage: s,
		broker:  b,
		houser:  h,
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
