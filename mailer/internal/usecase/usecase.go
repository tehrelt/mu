package usecase

import (
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/mailer/internal/config"
	"github.com/tehrelt/mu/mailer/internal/storage/mailer"
	"github.com/tehrelt/mu/mailer/pkg/pb/notificationpb"
)

type UseCase struct {
	cfg    *config.Config
	client notificationpb.NotificationServiceClient
	logger *slog.Logger
	mailer *mailer.Mailer
}

func New(
	cfg *config.Config,
	client notificationpb.NotificationServiceClient,
	mailer *mailer.Mailer,
) *UseCase {
	return &UseCase{
		cfg:    cfg,
		client: client,
		logger: slog.With(sl.Module("UseCase")),
		mailer: mailer,
	}
}
