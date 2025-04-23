package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/housing-service/internal/config"
	"github.com/tehrelt/mu/housing-service/internal/dto"
	"github.com/tehrelt/mu/housing-service/internal/storage"
	"github.com/tehrelt/mu/housing-service/internal/storage/pg/housestorage"
	"github.com/tehrelt/mu/housing-service/pkg/pb/housepb"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	cfg     *config.Config
	storage *housestorage.HouseStorage

	housepb.UnimplementedHouseServiceServer
}

func (s *Server) Create(ctx context.Context, in *housepb.CreateRequest) (*housepb.CreateResponse, error) {
	id, err := s.storage.Create(ctx, &dto.CreateHouse{
		Address:     in.Address,
		ResidentQty: int(in.ResidentsQty),
		RoomsQty:    int(in.RoomsQty),
	})
	if err != nil {
		slog.Error("failed to create house", sl.Err(err))
		return nil, err
	}

	return &housepb.CreateResponse{
		Id: id.String(),
	}, nil
}

func (s *Server) House(ctx context.Context, in *housepb.HouseRequest) (*housepb.HouseResponse, error) {

	id, err := uuid.Parse(in.HouseId)
	if err != nil {
		slog.Error("failed to parse id", sl.Err(err))
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	house, err := s.storage.Find(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrHouseNotFound) {
			return nil, status.Error(codes.NotFound, "house not found")
		}
		slog.Error("failed to find house", sl.Err(err))
		return nil, err
	}

	res := &housepb.HouseResponse{
		House: &housepb.House{
			Id:           string(house.Id.String()),
			Address:      house.Address,
			RoomsQty:     int64(house.RoomsQuantity),
			ResidentsQty: int64(house.ResidentsQuantity),
			CreatedAt:    house.CreatedAt.Unix(),
			ConnectedServices: lo.Map(house.ConnectedServices, func(item uuid.UUID, _ int) string {
				return item.String()
			}),
		},
	}

	if house.UpdatedAt != nil {
		res.House.UpdatedAt = house.UpdatedAt.Unix()
	}

	return res, nil
}

func New(cfg *config.Config, s *housestorage.HouseStorage) *Server {
	return &Server{
		cfg:     cfg,
		storage: s,
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

	housepb.RegisterHouseServiceServer(server, s)

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
