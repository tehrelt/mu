package authservice

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/lib/jwt"
	"github.com/tehrelt/mu/auth-service/internal/services"
	"github.com/tehrelt/mu/auth-service/internal/storage"
)

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.TokenPair, error) {

	fn := "authservice.Refresh"
	log := s.logger.With(sl.Method(fn))

	log.Debug("refreshing token", slog.String("refresh_token", refreshToken))

	claims, err := s.jwtClient.Verify(refreshToken, jwt.RefreshToken)
	if err != nil {
		return nil, err
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil {
		slog.Error("failed to parse user id from jwt payload", slog.String("claims.Id", claims.Id))
		return nil, err
	}

	if err := s.sessions.Check(ctx, userId, refreshToken); err != nil {
		if errors.Is(err, storage.ErrSessionInvalid) {
			return nil, services.ErrInvalidSession
		}

		log.Error("unexpected error", sl.Err(err))
		return nil, err
	}

	tokens, err := s.createTokens(&dto.UserClaims{Id: userId.String()})
	if err != nil {
		log.Error("cannot generate jwt", sl.Err(err))
		return nil, err
	}

	if err := s.sessions.Save(ctx, userId, tokens.RefreshToken); err != nil {
		log.Error("cannot save session", sl.Err(err))
		return nil, err
	}

	return tokens, nil
}
