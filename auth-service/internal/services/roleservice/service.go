package roleservice

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu/auth-service/internal/dto"
)

type RoleUpdated interface {
	Add(ctx context.Context, in *dto.UserRoles) error
	Remove(ctx context.Context, in *dto.UserRoles) error
}

type Service struct {
	updater RoleUpdated
}

func New(updater RoleUpdated) *Service {
	return &Service{
		updater: updater,
	}
}

func (s *Service) Add(ctx context.Context, in *dto.UserRoles) error {
	fn := "roleservice.Add"
	log := slog.With(slog.String("fn", fn))

	log.Info("Adding roles", slog.Any("roles", in))
	return s.updater.Add(ctx, in)
}

func (s *Service) Remove(ctx context.Context, in *dto.UserRoles) error {
	fn := "roleservice.Remove"
	log := slog.With(slog.String("fn", fn))

	log.Info("Removing roles", slog.Any("roles", in))
	return s.updater.Remove(ctx, in)
}
