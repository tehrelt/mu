package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tehrelt/mu/housing-service/pkg/prettyslog"
	"github.com/tehrelt/mu/housing-service/pkg/sl"
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

type Config struct {
	Env      Env `env:"ENV"`
	App      App
	Grpc     Grpc
	Postgres Postgres

	Jaeger struct {
		Endpoint string `env:"JAEGER_ENDPOINT"`
	}

	RateService struct {
		Host string `env:"RATE_SERVICE_HOST"`
		Port int    `env:"RATE_SERVICE_PORT"`
	}

	Amqp struct {
		Host string `env:"AMQP_HOST"`
		Port int    `env:"AMQP_PORT"`
		User string `env:"AMQP_USER"`
		Pass string `env:"AMQP_PASS"`
	}

	ConnectServiceExchange struct {
		Exchange string `env:"EXCHANGE"`
		Queue    string `env-default:"housing_service.connect_service"`
		Routing  string `env-default:"create_cabinet"`
	} `env-prefix:"RMQ_CONNECT_SERVICE_"`
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
