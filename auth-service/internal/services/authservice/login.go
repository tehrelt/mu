package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/domain"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/domain/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/domain/entity"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/storage"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

func (a *AuthService) Login(ctx context.Context, req *dto.Login) (*dto.Tokens, error) {

	fn := "authservice.Login"
	log := a.logger.With(sl.Method(fn))

	log.Debug("logging in", slog.Any("req", req))

	user, err := a.userProvider.Find(ctx, req.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Debug("user not found", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", fn, domain.ErrUserNotFound)
		}

		log.Error("unexpected error on found user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if err := a.comparePassword(user.HashedPassword, req.Password); err != nil {
		log.Error("password not match", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", fn, domain.ErrIncorrectPassword)
	}

	tokens, err := a.generateJwtPair(&entity.UserClaims{Id: user.Id, Email: user.Email})
	if err != nil {
		log.Error("cannot generate jwt", sl.Err(err))
		return nil, err
	}

	if err := a.sessions.Save(ctx, user.Id, tokens.RefreshToken); err != nil {
		log.Error("cannot save session", sl.Err(err))
		return nil, err
	}

	return tokens, nil
}
