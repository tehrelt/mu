package usecase

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/telegram-bot/internal/events"
	"github.com/tehrelt/mu/telegram-bot/internal/usecase/formatters"
	"go.opentelemetry.io/otel"
	"gopkg.in/telebot.v4"
)

func (uc *UseCase) SendNotification(ctx context.Context, event events.Event) error {

	fn := "SendNotification"
	log := slog.With(sl.Method(fn))
	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	defer span.End()

	chatId, err := strconv.ParseInt(event.Header().Settings.TelegramChatId, 10, 64)
	if err != nil {
		return err
	}

	formatter := formatters.New(uc.cfg)

	msg := formatter.Format(event)

	log.Info("sending message", slog.String("message", msg))
	_, err = uc.bot.Send(
		&telebot.Chat{
			ID: chatId,
		},
		msg,
		telebot.ParseMode(telebot.ModeMarkdownV2),
	)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to send message", slog.String("error", err.Error()))
	}

	return nil
}
