package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/tehrelt/mu-lib/sl"
)

type App struct {
	Name    string `env:"APP_NAME"`
	Version string `env:"APP_VERSION"`
}

type Env string

const (
	EnvLocal Env = "local"
	EnvDev   Env = "dev"
	EnvProd  Env = "prod"
)

type AMQP struct {
	Host  string `env:"HOST"`
	Port  int    `env:"PORT"`
	User  string `env:"USER"`
	Pass  string `env:"PASS"`
	Vhost string `env:"VHOST"`
}

func (a *AMQP) ConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", a.User, a.Pass, a.Host, a.Port, a.Vhost)
}

type Config struct {
	Env Env `env:"ENV"`
	App App

	AMQP AMQP `env-prefix:"AMQP_"`

	BotToken string `env:"BOT_TOKEN"`

	NotificationService struct {
		Host string `env:"NOTIFICATION_SERVICE_HOST"`
		Port int    `env:"NOTIFICATION_SERVICE_PORT"`
	}

	NotificationSendExchange string `env:"RMQ_NOTIFICATION_SEND_EXCHANGE"`

	Jaeger struct {
		Endpoint string `env:"JAEGER_ENDPOINT"`
	}

	SMTP struct {
		Host     string `env:"SMTP_HOST"`
		Port     int    `env:"SMTP_PORT"`
		Email    string `env:"SMTP_EMAIL"`
		Password string `env:"SMTP_PASSWORD"`
	}

	FrontendLink string `env:"FRONTEND_LINK"`
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
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	slog.SetDefault(log)
}
