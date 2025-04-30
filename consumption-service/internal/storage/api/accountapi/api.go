package accountapi

import (
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/consumption-service/internal/config"
	"github.com/tehrelt/mu/consumption-service/internal/usecase"
	"github.com/tehrelt/mu/consumption-service/pkg/pb/accountpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ usecase.AccountProvider = (*Api)(nil)

type Api struct {
	client accountpb.AccountServiceClient
	logger *slog.Logger
}

func New(config *config.Config) (*Api, error) {
	conn, err := grpc.NewClient(
		config.AccountService.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &Api{
		client: accountpb.NewAccountServiceClient(conn),
		logger: slog.With(sl.Module("accountapi")),
	}, nil
}
