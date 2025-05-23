// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"context"
	"fmt"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/gateway/internal/config"
	http2 "github.com/tehrelt/mu/gateway/internal/transport/http/admin"
	"github.com/tehrelt/mu/gateway/internal/transport/http/public"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/billingpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/consumptionpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/notificationpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/ratepb"
	"github.com/tehrelt/mu/gateway/pkg/pb/registerpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/ticketpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/userpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

import (
	_ "github.com/jackc/pgx/stdlib"
)

// Injectors from wire.go:

func New(ctx context.Context) (*App, func(), error) {
	configConfig := config.New()
	authServiceClient, err := _authpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	registerServiceClient, err := _regpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	accountServiceClient, err := _accountpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	rateServiceClient, err := _ratepb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	userServiceClient, err := _userpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	billingServiceClient, err := _billingpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	ticketServiceClient, err := _ticketpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	consumptionServiceClient, err := _consumptionpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	notificationServiceClient, err := _notificationpb(configConfig)
	if err != nil {
		return nil, nil, err
	}
	server := http.New(configConfig, authServiceClient, registerServiceClient, accountServiceClient, rateServiceClient, userServiceClient, billingServiceClient, ticketServiceClient, consumptionServiceClient, notificationServiceClient)
	httpServer := http2.New(configConfig, authServiceClient, registerServiceClient, accountServiceClient, rateServiceClient, userServiceClient, billingServiceClient, ticketServiceClient)
	v, err := _servers(server, httpServer)
	if err != nil {
		return nil, nil, err
	}
	tracer, err := _tracer(ctx, configConfig)
	if err != nil {
		return nil, nil, err
	}
	app := newApp(v, tracer)
	return app, func() {
	}, nil
}

// wire.go:

func _notificationpb(cfg *config.Config) (notificationpb.NotificationServiceClient, error) {
	host := cfg.NotificationService.Host
	port := cfg.NotificationService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return notificationpb.NewNotificationServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(notificationpb.NotificationServiceClient), nil
}

func _consumptionpb(cfg *config.Config) (consumptionpb.ConsumptionServiceClient, error) {
	host := cfg.ConsumptionService.Host
	port := cfg.ConsumptionService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return consumptionpb.NewConsumptionServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(consumptionpb.ConsumptionServiceClient), nil
}

func _ticketpb(cfg *config.Config) (ticketpb.TicketServiceClient, error) {
	host := cfg.TicketService.Host
	port := cfg.TicketService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return ticketpb.NewTicketServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(ticketpb.TicketServiceClient), nil
}

func _authpb(cfg *config.Config) (authpb.AuthServiceClient, error) {
	host := cfg.AuthService.Host
	port := cfg.AuthService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return authpb.NewAuthServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(authpb.AuthServiceClient), nil
}

func _billingpb(cfg *config.Config) (billingpb.BillingServiceClient, error) {
	host := cfg.BillingService.Host
	port := cfg.BillingService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return billingpb.NewBillingServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(billingpb.BillingServiceClient), nil
}

func _userpb(cfg *config.Config) (userpb.UserServiceClient, error) {
	host := cfg.UserService.Host
	port := cfg.UserService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return userpb.NewUserServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(userpb.UserServiceClient), nil
}

func _accountpb(cfg *config.Config) (accountpb.AccountServiceClient, error) {
	host := cfg.AccountService.Host
	port := cfg.AccountService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return accountpb.NewAccountServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(accountpb.AccountServiceClient), nil
}

func _regpb(cfg *config.Config) (registerpb.RegisterServiceClient, error) {
	host := cfg.RegisterService.Host
	port := cfg.RegisterService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return registerpb.NewRegisterServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(registerpb.RegisterServiceClient), nil
}

func _ratepb(cfg *config.Config) (ratepb.RateServiceClient, error) {
	host := cfg.RateService.Host
	port := cfg.RateService.Port

	client, err := create_grpc_client(host, port, func(cc grpc.ClientConnInterface) any {
		return ratepb.NewRateServiceClient(cc)
	})
	if err != nil {
		return nil, err
	}

	return client.(ratepb.RateServiceClient), nil
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	return tracer.SetupTracer(ctx, cfg.Jaeger.Endpoint, cfg.App.Name)
}

func create_grpc_client(host string, port int, fn func(grpc.ClientConnInterface) any) (any, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClient(
		addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()), grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return fn(conn), nil
}

func _servers(public *http.Server, admin *http2.Server) ([]Server, error) {
	servers := []Server{public, admin}
	return servers, nil
}
