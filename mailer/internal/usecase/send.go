package usecase

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/mailer/internal/dto"
	"github.com/tehrelt/mu/mailer/internal/events"
	"github.com/tehrelt/mu/mailer/internal/usecase/formatters"
	"go.opentelemetry.io/otel"
)

func (uc *UseCase) SendNotification(ctx context.Context, event events.Event) (err error) {

	fn := "SendNotification"
	log := slog.With(sl.Method(fn))
	ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fn)
	defer func() {
		if err != nil {
			span.RecordError(err)
		}

		span.End()
	}()

	to := event.Header().Settings.Email

	formatter := formatters.New(uc.cfg)
	header, msg := formatter.Format(event)

	log.Info("sending message", slog.String("to", to), slog.String("message", msg))

	if err := uc.mailer.Send(ctx, &dto.Message{
		To:      []string{to},
		Subject: header,
		Body:    msg,
	}); err != nil {
		log.Error("failed to send message", slog.String("to", to), slog.String("message", msg), slog.String("error", err.Error()))
		return err
	}

	return nil
}
