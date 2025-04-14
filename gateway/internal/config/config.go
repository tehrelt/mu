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

type Config struct {
	Env Env `env:"ENV"`
	App App

	Http struct {
		Host string `env:"HTTP_HOST"`
		Port int    `env:"HTTP_PORT"`
	}

	RegisterService struct {
		Host string `env:"REGISTER_SERVICE_HOST"`
		Port int    `env:"REGISTER_SERVICE_PORT"`
	}

	BillingService struct {
		Host string `env:"BILLING_SERVICE_HOST"`
		Port int    `env:"BILLING_SERVICE_PORT"`
	}

	RateService struct {
		Host string `env:"RATE_SERVICE_HOST"`
		Port int    `env:"RATE_SERVICE_PORT"`
	}

	AccountService struct {
		Host string `env:"ACCOUNT_SERVICE_HOST"`
		Port int    `env:"ACCOUNT_SERVICE_PORT"`
	}

	UserService struct {
		Host string `env:"USER_SERVICE_HOST"`
		Port int    `env:"USER_SERVICE_PORT"`
	}

	AuthService struct {
		Host string `env:"AUTH_SERVICE_HOST"`
		Port int    `env:"AUTH_SERVICE_PORT"`
	}

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
