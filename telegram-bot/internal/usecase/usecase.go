package usecase

import (
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/telegram-bot/pkg/pb/notificationpb"
	"gopkg.in/telebot.v4"
)

type UseCase struct {
	client notificationpb.NotificationServiceClient
	logger *slog.Logger
	bot    *telebot.Bot
}

func New(client notificationpb.NotificationServiceClient, bot *telebot.Bot) *UseCase {
	return &UseCase{
		client: client,
		logger: slog.With(sl.Module("UseCase")),
		bot:    bot,
	}
}
