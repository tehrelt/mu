package mailer

import (
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/mailer/internal/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	cfg *config.Config
	l   *slog.Logger
}

func New(cfg *config.Config) *Mailer {
	return &Mailer{
		cfg: cfg,
		l:   slog.With(sl.Module("mailer")),
	}
}

func (m *Mailer) dial() (*gomail.Dialer, error) {
	fn := "dial"
	log := m.l.With(sl.Method(fn))

	log.Debug("dialing smtp")
	dialer := gomail.NewDialer(
		m.cfg.SMTP.Host,
		m.cfg.SMTP.Port,
		m.cfg.SMTP.Email,
		m.cfg.SMTP.Password,
	)

	log.Info("dialing smtp")
	conn, err := dialer.Dial()
	if err != nil {
		log.Error("failed to dial smtp", sl.Err(err))
		return nil, err
	}
	defer conn.Close()

	return dialer, nil
}
