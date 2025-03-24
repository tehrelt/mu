package authservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/lib/jwt"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/services"
)

func (s *AuthService) Authorize(ctx context.Context, token string, roles ...models.Role) (*dto.UserClaims, error) {
	claims, err := s.jwtClient.Verify(token, jwt.AccessToken)
	if err != nil {
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

	for _, userRole := range userRoles {
		for _, targetRole := range roles {
			if userRole == targetRole {
				return claims, nil
			}
		}
	}

	return nil, services.ErrForbidden
}
