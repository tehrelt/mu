package authservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/storage"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterUser) (*dto.TokenPair, error) {

	log := slog.With(sl.Method("authservice.Register"))

	candidate, err := s.userProvider.UserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, storage.ErrUserNotFound) {
		log.Error("get user by email error", sl.Err(err))
		return nil, err
	}
	if candidate != nil {
		log.Error("user already exists")
		return nil, fmt.Errorf("user already exists")
	}

	user := &dto.CreateUser{
		Fio:          req.Fio,
		PersonalData: req.PersonalData,
		Email:        req.Email,
	}

	hashedPassword, err := s.hash(req.Password)
	if err != nil {
		log.Error("hash error", sl.Err(err))
		return nil, err
	}

	userId, err := s.userSaver.Create(ctx, user)
	if err != nil {
		log.Error("create user error", sl.Err(err))
		return nil, err
	}

	creds := &models.Credentials{
		UserId:         userId,
		HashedPassword: hashedPassword,
		Roles:          req.Roles,
	}

	if err := s.credentialSaver.Save(ctx, creds); err != nil {
		log.Error("save credentials error", sl.Err(err))
		return nil, err
	}

	tokens, err := s.createTokens(&dto.UserClaims{Id: userId.String()})
	if err != nil {
		log.Error("create tokens error", sl.Err(err))
		return nil, err
	}

	if err := s.sessions.Save(ctx, userId, tokens.RefreshToken); err != nil {
		log.Error("save session error", sl.Err(err))
		return nil, err
	}

	return tokens, nil
}
