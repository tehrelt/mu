package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/consumption-service/pkg/prettyslog"
)

type App struct {
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION"`
}

type ServerConfig struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port int    `env:"PORT"`
}

type Database struct {
	Proto    string `env:"PROTO" env-default:"postgresql"`
	Host     string `env:"HOST" env-required:"true"`
	Port     int    `env:"PORT" env-required:"true"`
	User     string `env:"USER"`
	Pass     string `env:"PASS"`
	Database string `env:"DATABASE" env-required:"true"`
}

type ApiConfig struct {
	Host string `env:"HOST" env-default:"0.0.0.0"`
	Port int    `env:"PORT"`
}

func (ac *ApiConfig) Address() string {
	return fmt.Sprintf("%s:%d", ac.Host, ac.Port)
}

func (m *Database) ConnectionString() string {
	if m.User == "" || m.Pass == "" {
		return fmt.Sprintf("%s://%s:%d", m.Proto, m.Host, m.Port)
	}

	return fmt.Sprintf("%s://%s:%s@%s:%d/%s", m.Proto, m.User, m.Pass, m.Host, m.Port, m.Database)
}

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type Config struct {
	Env  Env `env:"ENV"`
	App  App
	Grpc ServerConfig `env-prefix:"GRPC_"`

	BillingService ApiConfig `env-prefix:"BILLING_SERVICE_"`
	RateService    ApiConfig `env-prefix:"RATE_SERVICE_"`
	AccountService ApiConfig `env-prefix:"ACCOUNT_SERVICE_"`

	Postgres Database `env-prefix:"PG_"`

	Jaeger struct {
		Endpoint string `env:"JAEGER_ENDPOINT"`
	}

	Amqp struct {
		Host string `env:"AMQP_HOST"`
		Port int    `env:"AMQP_PORT"`
		User string `env:"AMQP_USER"`
		Pass string `env:"AMQP_PASS"`
	}

	ServiceConnectedQueue struct {
		Exchange string `env:"SERVICE_CONNECTED_QUEUE_EXCHANGE"`
		Routing  string `env:"SERVICE_CONNECTED_QUEUE_ROUTING"`
	}
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
