package cron

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/register-service/pkg/pb/registerpb"
)

func (s *Cron) reg(ctx context.Context) error {

	has, err := s.regservice.HasAdmin(ctx)
	if err != nil {
		slog.Error("failed to check admin", sl.Err(err))
		return err
	}

	if has {
		slog.Info("admin found, skip register default")
		return nil
	}

	slog.Info("no admin found, creating admin")
	if _, err = s.regservice.Register(ctx, &registerpb.RegisterRequest{
		User: &registerpb.User{
			LastName:       "root",
			FirstName:      "",
			MiddleName:     "",
			Email:          s.cfg.DefaultAdmin.Email,
			Phone:          "70000000000",
			PassportNumber: 0,
			PassportSeries: 0,
			Snils:          "000-0000-000",
			Password:       s.cfg.DefaultAdmin.Password,
			Roles:          []registerpb.Role{registerpb.Role_ROLE_ADMIN},
		},
	}); err != nil {
		slog.Error("failed to create admin", sl.Err(err))
		return err
	}

	return nil
}
