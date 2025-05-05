package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/notification-service/pkg/prettyslog"
)

type App struct {
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION"`
}

type Grpc struct {
	Host string `env:"GRPC_HOST" env-default:"0.0.0.0"`
	Port int    `env:"GRPC_PORT" env-default:"8080"`
}

type Database struct {
	Host string `env:"HOST" env-required:"true"`
	Port int    `env:"PORT" env-required:"true"`
	User string `env:"USER"`
	Pass string `env:"PASS"`
	Name string `env:"NAME"`
}

func (m *Database) ConnectionString() string {
	return fmt.Sprintf("%s:%s@%s:%d/%s", m.User, m.Pass, m.Host, m.Port, m.Name)
}

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type ExternalServiceConfig struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

func (e ExternalServiceConfig) Address() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

type Config struct {
	Env      Env `env:"ENV"`
	App      App
	Grpc     Grpc
	Postgres Database `env-prefix:"PG_"`
	Redis    Database `env-prefix:"REDIS_"`

	Jaeger struct {
		Endpoint string `env:"JAEGER_ENDPOINT"`
	}

	TicketService  ExternalServiceConfig `env-prefix:"TICKET_SERVICE_"`
	UserService    ExternalServiceConfig `env-prefix:"USER_SERVICE_"`
	AccountService ExternalServiceConfig `env-prefix:"ACCOUNT_SERVICE_"`

	Amqp struct {
		Host string `env:"AMQP_HOST"`
		Port int    `env:"AMQP_PORT"`
		User string `env:"AMQP_USER"`
		Pass string `env:"AMQP_PASS"`
	}

	TicketStatusChangedExchange struct {
		Exchange string `env:"EXCHANGE"`
	} `env-prefix:"RMQ_TICKET_STATUS_CHANGED_"`

	BalanceChangedExchange string `env:"RMQ_BALANCE_CHANGED_EXCHANGE"`

	NotificationSendExchange struct {
		Exchange string `env:"EXCHANGE"`
	} `env-prefix:"RMQ_NOTIFICATION_SEND_"`
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

	interceptors.SetDebug(false)

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
