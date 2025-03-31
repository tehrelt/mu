package authservice

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/models"
)

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterUser) (*dto.TokenPair, error) {

	log := slog.With(sl.Method("authservice.Register"))

	hashedPassword, err := s.hash(req.Password)
	if err != nil {
		log.Error("hash error", sl.Err(err))
		return nil, err
	}

	creds := &models.Credentials{
		UserId:         req.UserId,
		HashedPassword: hashedPassword,
		Roles:          req.Roles,
	}

	if err := s.credentialSaver.Save(ctx, creds); err != nil {
		log.Error("save credentials error", sl.Err(err))
		return nil, err
	}

	tokens, err := s.createTokens(&dto.UserClaims{Id: req.UserId.String()})
	if err != nil {
		log.Error("create tokens error", sl.Err(err))
		return nil, err
	}

	if err := s.sessions.Save(ctx, req.UserId, tokens.RefreshToken); err != nil {
		log.Error("save session error", sl.Err(err))
		return nil, err
	}

	return tokens, nil
}
