package profileservice

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/lib/jwt"
	"github.com/tehrelt/mu/auth-service/internal/services"
)

func (s *ProfileService) Profile(ctx context.Context, token string) (*dto.Profile, error) {

	log := slog.With(sl.Method("profileservice.Profile"))

	payload, err := s.jc.Verify(token, jwt.AccessToken)
	if err != nil {
		log.Error("failed to jwt verify", sl.Err(err))
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, services.ErrTokenExpired
		}
		return nil, err
	}

	userId, err := uuid.Parse(payload.Id)
	if err != nil {
		log.Error("failed to parse uuid", sl.Err(err))
		return nil, err
	}

	user, err := s.userProvider.UserById(ctx, userId)
	if err != nil {
		log.Error("failed to get user by id", sl.Err(err))
		return nil, err
	}

	roles, err := s.roleProvider.Roles(ctx, userId)
	if err != nil {
		log.Error("failed to get user roles", sl.Err(err))
		return nil, err
	}

	return &dto.Profile{
		Id:    user.Id,
		Email: user.Email,
		Roles: roles,
		Fio: dto.Fio{
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName,
		},
	}, nil
}
