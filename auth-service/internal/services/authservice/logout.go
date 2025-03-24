package authservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

func (s *AuthService) Logout(ctx context.Context, userId uuid.UUID) error {

	log := slog.With(sl.Method("authservice.Logout"))

	if err := s.sessions.Delete(ctx, userId); err != nil {
		log.Error("failed to delete session", sl.Err(err))
		return fmt.Errorf("failed delete session: %w", err)
	}

	return nil
}
