package authservice

import (
	"context"
	"fmt"
	"log/slog"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/domain/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/domain/entity"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

func (s *AuthService) AddRoles(ctx context.Context, slug string, roles []entity.Role) error {

	fn := "authservice.AddRoles"
	log := s.logger.With(sl.Method(fn), slog.String("userSlug", slug), slog.Any("roles", roles))

	if err := s.roles.Add(ctx, &dto.AddRoles{
		UserId: slug,
		Roles:  roles,
	}); err != nil {
		log.Error("failed to add roles", sl.Err(err))
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}
