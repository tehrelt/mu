package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/account-service/pkg/prettyslog"
)

type App struct {
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION"`
}

type Grpc struct {
	Host string `env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port int    `env:"GRPC_PORT" env-default:"8080"`
}

type Postgres struct {
	Host string `env:"PG_HOST" env-required:"true"`
	Port int    `env:"PG_PORT" env-required:"true"`
	User string `env:"PG_USER" env-required:"true"`
	Pass string `env:"PG_PASS" env-required:"true"`
	Name string `env:"PG_NAME" env-required:"true"`
}

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type QueueConfig struct {
	Exchange string `env:"EXCHANGE"`
	Routing  string `env:"ROUTING_KEY"`
}

type ExternalServiceConfig struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

type Config struct {
	Env      Env `env:"ENV"`
	App      App
	Grpc     Grpc
	Postgres Postgres

	Jaeger struct {
		Endpoint string `env:"JAEGER_ENDPOINT"`
	}

	RateService    ExternalServiceConfig `env-prefix:"RATE_SERVICE_"`
	UserService    ExternalServiceConfig `env-prefix:"USER_SERVICE_"`
	HouseService   ExternalServiceConfig `env-prefix:"HOUSE_SERVICE_"`
	BillingService ExternalServiceConfig `env-prefix:"BILLING_SERVICE_"`
	TicketService  ExternalServiceConfig `env-prefix:"TICKET_SERVICE_"`

	Amqp struct {
		Host string `env:"AMQP_HOST"`
		Port int    `env:"AMQP_PORT"`
		User string `env:"AMQP_USER"`
		Pass string `env:"AMQP_PASS"`
	}

	PaymentStatusChanged struct {
		Exchange string `env:"PAYMENT_STATUS_CHANGED_EXCHANGE"`
		Routing  string `env:"PAYMENT_STATUS_CHANGED_ROUTING"`
	}

	BalanceChanged struct {
		Exchange string `env:"BALANCE_CHANGED_EXCHANGE"`
		Routing  string `env:"BALANCE_CHANGED_ROUTING"`
	}

	TicketStatusChanged struct {
		Exchange            string `env:"EXCHANGE"`
		NewAccountRoute     string `env:"NEW_ACCOUNT_ROUTE"`
		ConnectServiceRoute string `env:"CONNECT_SERVICE_ROUTE"`
	} `env-prefix:"TICKET_STATUS_CHANGED_QUEUE_"`
}

func New() *Config {
	config := new(Config)

	if err := cleanenv.ReadEnv(config); err != nil {
		slog.Error("error when reading env", sl.Err(err))
		header := fmt.Sprintf("%s - %s", os.Getenv("APP_NAME"), os.Getenv("APP_VERSION"))

		usage := cleanenv.FUsage(os.Stdout, config, &header)
		usage()

		os.Exit(-1)
	}

	setupLogger(config)

	slog.Debug("config", slog.Any("c", config))
	return config
}

func setupLogger(cfg *Config) {
	var log *slog.Logger

	switch cfg.Env {
	case EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case EnvDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(prettyslog.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	slog.SetDefault(log)
}
