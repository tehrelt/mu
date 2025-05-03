package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/gateway/pkg/prettyslog"
)

type App struct {
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION"`
}

type Http struct {
	Host string `env:"HTTP_HOST" env-required:"true"`
	Port int    `env:"HTTP_PORT" env-required:"true"`
}

type Cors struct {
	AllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" env-required:"true"`
}

func (c *Cors) Split() []string {
	return strings.Split(c.AllowedOrigins, ",")
}

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type HttpConfig struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

type ExternalServiceConfig struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

func (svc *ExternalServiceConfig) Address() string {
	return fmt.Sprintf("%s:%d", svc.Host, svc.Port)
}

type Config struct {
	Env Env `env:"ENV"`
	App App

	Cors Cors

	PublicHttpApi HttpConfig `env-prefix:"PUBLIC_HTTP_"`
	AdminHttpApi  HttpConfig `env-prefix:"ADMIN_HTTP_"`

	RegisterService     ExternalServiceConfig `env-prefix:"REGISTER_SERVICE_"`
	BillingService      ExternalServiceConfig `env-prefix:"BILLING_SERVICE_"`
	RateService         ExternalServiceConfig `env-prefix:"RATE_SERVICE_"`
	AccountService      ExternalServiceConfig `env-prefix:"ACCOUNT_SERVICE_"`
	UserService         ExternalServiceConfig `env-prefix:"USER_SERVICE_"`
	AuthService         ExternalServiceConfig `env-prefix:"AUTH_SERVICE_"`
	TicketService       ExternalServiceConfig `env-prefix:"TICKET_SERVICE_"`
	ConsumptionService  ExternalServiceConfig `env-prefix:"CONSUMPTION_SERVICE_"`
	NotificationService ExternalServiceConfig `env-prefix:"NOTIFICATION_SERVICE_"`

	Jaeger struct {
		Endpoint string `env:"JAEGER_ENDPOINT"`
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

	if config.Env == EnvLocal {
		// interceptors.SetDebug(true)
	}

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
