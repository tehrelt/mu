package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/telegram-bot/internal/dto"
	"github.com/tehrelt/mu/telegram-bot/pkg/pb/notificationpb"
)

func (u *UseCase) Link(ctx context.Context, in *dto.LinkUser) error {

	fn := "Link"
	log := u.logger.With(sl.Method(fn))

	log.Info("linking telegram", slog.Any("in", in))
	if _, err := u.client.LinkTelegram(ctx, &notificationpb.LinkTelegramRequest{
		UserId: in.UserId,
		ChatId: fmt.Sprintf("%d", in.ChatId),
		Otp:    in.Code,
	}); err != nil {
		log.Error("failed linking", sl.Err(err))
		return err
	}

	return nil
}
