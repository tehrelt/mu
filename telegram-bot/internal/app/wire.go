//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"errors"

	"github.com/google/wire"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/telegram-bot/internal/config"
	"github.com/tehrelt/mu/telegram-bot/internal/transport/tg"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/telebot.v4"
)

func NewApp(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		tg.New,

		_bot,
		_tracer,
		config.New,
	))
}

func _tracer(ctx context.Context, cfg *config.Config) (trace.Tracer, error) {
	return tracer.SetupTracer(ctx, cfg.Jaeger.Endpoint, cfg.App.Name)
}

func _bot(ctx context.Context, cfg *config.Config) (*telebot.Bot, error) {

	token := cfg.BotToken

	if token == "" {
		return nil, errors.New("bot token is empty")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return bot, nil
}
