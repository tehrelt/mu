package authservice

import (
	"context"
	"errors"
	"log/slog"

	"slices"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/lib/jwt"
	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/internal/services"
)

func (s *AuthService) Authorize(ctx context.Context, token string, roles ...models.Role) (*dto.UserClaims, error) {
	log := s.logger.With(sl.Method("Authorize"))
	claims, err := s.jwtClient.Verify(token, jwt.AccessToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, services.ErrTokenExpired
		}
		return nil, err
	}

	userid, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return claims, nil
	}

	userRoles, err := s.roles.Roles(ctx, userid)
	if err != nil {
		return nil, err
	}

	log.Debug("checking role", slog.Any("user_roles", userRoles), slog.Any("allowed_roles", roles))
	for _, userRole := range userRoles {
		if slices.Contains(roles, userRole) {
			return claims, nil
		}
	}

	return nil, services.ErrForbidden
}
