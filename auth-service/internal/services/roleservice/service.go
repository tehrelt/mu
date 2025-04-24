package roleservice

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/models"
)

type RoleStorage interface {
	Add(ctx context.Context, in *dto.UserRoles) error
	Remove(ctx context.Context, in *dto.UserRoles) error
	Count(ctx context.Context, role ...models.Role) (int64, error)
}

type Service struct {
	storage RoleStorage
}

func New(updater RoleStorage) *Service {
	return &Service{
		storage: updater,
	}
}

func (s *Service) Count(ctx context.Context, role ...models.Role) (int64, error) {
	fn := "roleservice.Count"
	log := slog.With(slog.String("fn", fn))

	log.Info("Counting roles", slog.Any("roles", role))
	return s.storage.Count(ctx, role...)
}

func (s *Service) Add(ctx context.Context, in *dto.UserRoles) error {
	fn := "roleservice.Add"
	log := slog.With(slog.String("fn", fn))

	log.Info("Adding roles", slog.Any("roles", in))
	return s.storage.Add(ctx, in)
}

func (s *Service) Remove(ctx context.Context, in *dto.UserRoles) error {
	fn := "roleservice.Remove"
	log := slog.With(slog.String("fn", fn))

	log.Info("Removing roles", slog.Any("roles", in))
	return s.storage.Remove(ctx, in)
}
