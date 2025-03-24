package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/services"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/storage"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(ctx context.Context, req *dto.LoginUser) (*dto.TokenPair, error) {

	log := slog.With(sl.Method("authservice.Login"))

	candidate, err := s.userProvider.UserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, services.ErrInvalidCredentials
		}

		log.Error("failed provide user", slog.String("email", req.Email), sl.Err(err))
		return nil, fmt.Errorf("failed to provide user: %w", err)
	}

	creds, err := s.credentialsProvider.Credentials(ctx, candidate.Id)
	if err != nil {
		log.Error("failed to get credentials", slog.String("user_id", candidate.Id.String()), sl.Err(err))
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	if err := s.comparePassword(creds.HashedPassword, req.Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, services.ErrInvalidCredentials
		}

		log.Error("failed to compare password", slog.String("user_id", candidate.Id.String()), sl.Err(err))
		return nil, fmt.Errorf("failed to compare password: %w", err)
	}

	tokens, err := s.createTokens(&dto.UserClaims{Id: candidate.Id.String()})
	if err != nil {
		log.Error("failed to create tokens", slog.String("user_id", candidate.Id.String()), sl.Err(err))
		return nil, fmt.Errorf("failed to create tokens: %w", err)
	}

	if err := s.sessions.Save(ctx, candidate.Id, tokens.RefreshToken); err != nil {
		log.Error("failed to save session", slog.String("user_id", candidate.Id.String()), slog.String("refresh_token", tokens.RefreshToken), sl.Err(err))
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return tokens, nil
}
