package tg

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/tehrelt/mu/telegram-bot/internal/dto"
	"gopkg.in/telebot.v4"
)

func (b *Bot) startHandler() telebot.HandlerFunc {
	return func(ctx telebot.Context) error {

		sender := ctx.Sender()
		chat := ctx.Chat()
		args := ctx.Args()

		b.logger.Debug("start message",
			slog.String("sender", sender.FirstName),
			slog.Int64("chat_id", chat.ID),
			slog.Any("args", args),
		)

		if len(args) == 0 {
			return ctx.Send(fmt.Sprintf("Invalid start message missing token"))
		}

		payload := args[0]

		parts := strings.Split(payload, "_")
		if len(parts) != 2 {
			return ctx.Send(fmt.Sprintf("Invalid start message missing token"))
		}

		token := parts[0]
		userid := parts[1]

		user, err := b.uc.Link(context.Background(), &dto.LinkUser{
			ChatId: chat.ID,
			UserId: userid,
			Code:   token,
		})
		if err != nil {
			return ctx.Send(fmt.Sprintf("Failed to link user: %v", err))
		}

		return ctx.Send(fmt.Sprintf("Добро пожаловать, %s %s %s. \n\nЭто бот системы Мои услуги сюда будут приходить уведомления обновлении статуса созданных Вами заявок и уведомления операций над балансом.", user.Fio.Lastname, user.Fio.Firstname, user.Fio.Middlename))
	}
}
