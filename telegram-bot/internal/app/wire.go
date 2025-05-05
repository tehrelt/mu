//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/wire"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu-lib/tracer/interceptors"
	"github.com/tehrelt/mu/telegram-bot/internal/config"
	"github.com/tehrelt/mu/telegram-bot/internal/transport/amqp"
	"github.com/tehrelt/mu/telegram-bot/internal/transport/tg"
	"github.com/tehrelt/mu/telegram-bot/internal/usecase"
	"github.com/tehrelt/mu/telegram-bot/pkg/pb/notificationpb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/telebot.v4"
)

func NewApp(ctx context.Context) (*App, func(), error) {
	panic(wire.Build(
		newApp,
		tg.New,
		amqp.New,

		usecase.New,

		_notificationpb,

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

func _notificationpb(cfg *config.Config) (
	notificationpb.NotificationServiceClient,
	func(),
	error,
) {
	host := cfg.NotificationService.Host
	port := cfg.NotificationService.Port
	addr := fmt.Sprintf("%s:%d", host, port)

	client, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(interceptors.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, nil, err
	}

	return notificationpb.NewNotificationServiceClient(client), func() { client.Close() }, nil
}
