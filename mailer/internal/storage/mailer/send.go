package mailer

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/mailer/internal/dto"
	"gopkg.in/gomail.v2"
)

func (m *Mailer) Send(ctx context.Context, in *dto.Message) error {
	dialer, err := m.dial()
	if err != nil {
		return err
	}

	message := gomail.NewMessage()

	slog.Info("preparing messagew", slog.Any("dto", in))

	message.SetHeader("From", m.cfg.SMTP.Email)
	message.SetHeader("To", in.To...)
	message.SetHeader("Subject", in.Subject)
	message.SetBody("text/html", in.Body)

	if err := dialer.DialAndSend(message); err != nil {
		slog.Error("failed to send message", sl.Err(err))
		return err
	}

	slog.Info("message sent")

	return nil
}
