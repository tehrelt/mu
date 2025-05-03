package usecase

import (
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/telegram-bot/pkg/pb/notificationpb"
)

type UseCase struct {
	client notificationpb.NotificationServiceClient
	logger *slog.Logger
}

func New(client notificationpb.NotificationServiceClient) *UseCase {
	return &UseCase{
		client: client,
		logger: slog.With(sl.Module("UseCase")),
	}
}
