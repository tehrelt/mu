package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/telegram-bot/internal/dto"
	"github.com/tehrelt/mu/telegram-bot/pkg/pb/notificationpb"
	"github.com/tehrelt/mu/telegram-bot/pkg/pb/userpb"
)

func (u *UseCase) Link(ctx context.Context, in *dto.LinkUser) (*userpb.User, error) {

	fn := "Link"
	log := u.logger.With(sl.Method(fn))

	log.Debug("searching user for id", slog.String("userId", in.UserId))
	user, err := u.usersProvider.Find(ctx, &userpb.FindRequest{
		SearchBy: &userpb.FindRequest_Id{
			Id: in.UserId,
		},
	})
	if err != nil {
		log.Error("failed finding user", sl.Err(err))
		return nil, err
	}

	log.Info("linking telegram", slog.Any("in", in))
	if _, err := u.client.LinkTelegram(ctx, &notificationpb.LinkTelegramRequest{
		UserId: in.UserId,
		ChatId: fmt.Sprintf("%d", in.ChatId),
		Otp:    in.Code,
	}); err != nil {
		log.Error("failed linking", sl.Err(err))
		return nil, err
	}

	return user.User, nil
}
