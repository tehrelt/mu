package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/ticket-service/pkg/prettyslog"
)

type App struct {
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION"`
}

type Grpc struct {
	Host string `env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port int    `env:"GRPC_PORT" env-default:"8080"`
}

type Mongo struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     int    `env:"PORT" env-required:"true"`
	User     string `env:"USER"`
	Pass     string `env:"PASS"`
	Database string `env:"DATABASE" env-required:"true"`
}

func (m *Mongo) ConnectionString() string {
	if m.User == "" || m.Pass == "" {
		return fmt.Sprintf("mongodb://%s:%d", m.Host, m.Port)
	}

	return fmt.Sprintf("mongodb://%s:%s@%s:%d", m.User, m.Pass, m.Host, m.Port)
}

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Env   Env `env:"ENV"`
	App   App
	Grpc  Grpc
	Mongo Mongo `env-prefix:"MONGO_"`

	Jaeger struct {
		Endpoint string `env:"JAEGER_ENDPOINT"`
	}

	Amqp struct {
		Host string `env:"AMQP_HOST"`
		Port int    `env:"AMQP_PORT"`
		User string `env:"AMQP_USER"`
		Pass string `env:"AMQP_PASS"`
	}

	TicketStatusChangedQueue QueueConfig `env-prefix:"TICKET_STATUS_CHANGED_QUEUE_"`
}

type QueueConfig struct {
	RoutingKey string `env:"ROUTING_KEY"`
	Exchange   string `env:"EXCHANGE"`
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
