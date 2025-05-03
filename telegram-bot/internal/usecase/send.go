package usecase

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/tehrelt/mu/telegram-bot/internal/events"
	"github.com/tehrelt/mu/telegram-bot/internal/usecase/formatters"
	"gopkg.in/telebot.v4"
)

func (uc *UseCase) SendNotification(ctx context.Context, event events.Event) error {

	chatId, err := strconv.ParseInt(event.Header().Settings.TelegramChatId, 10, 64)
	if err != nil {
		return err
	}

	formatter := formatters.New(uc.cfg)

	msg := formatter.Format(event)

	slog.Info("sending message", slog.String("message", msg))
	_, err = uc.bot.Send(
		&telebot.Chat{
			ID: chatId,
		},
		msg,
		telebot.ParseMode(telebot.ModeMarkdownV2),
	)
	if err != nil {
		slog.Error("failed to send message", slog.String("error", err.Error()))
	}

	return nil
}
