package usecase

import (
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/telegram-bot/internal/config"
	"github.com/tehrelt/mu/telegram-bot/pkg/pb/notificationpb"
	"github.com/tehrelt/mu/telegram-bot/pkg/pb/userpb"
	"gopkg.in/telebot.v4"
)

type UseCase struct {
	cfg           *config.Config
	client        notificationpb.NotificationServiceClient
	usersProvider userpb.UserServiceClient
	logger        *slog.Logger
	bot           *telebot.Bot
}

func New(
	cfg *config.Config,
	client notificationpb.NotificationServiceClient,
	usersProvider userpb.UserServiceClient,
	bot *telebot.Bot,
) *UseCase {
	return &UseCase{
		cfg:           cfg,
		client:        client,
		usersProvider: usersProvider,
		logger:        slog.With(sl.Module("UseCase")),
		bot:           bot,
	}
}
