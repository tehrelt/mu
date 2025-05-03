package tg

import (
	"fmt"
	"log/slog"

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

		token := args[0]

		return ctx.Send(fmt.Sprintf("Welcome, %s!\nChat Id: %d\nToken: %s", sender.FirstName, chat.ID, token))
	}
}
