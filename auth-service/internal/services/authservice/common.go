package authservice

import (
	"github.com/tehrelt/mu/auth-service/internal/dto"
	"github.com/tehrelt/mu/auth-service/internal/lib/jwt"
	"golang.org/x/crypto/bcrypt"
)

func (a *AuthService) hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		a.cfg.Bcrypt.Cost,
	)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (a *AuthService) comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (a *AuthService) createTokens(claims *dto.UserClaims) (*dto.TokenPair, error) {
	accessToken, err := a.jwtClient.Sign(claims, jwt.AccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.jwtClient.Sign(claims, jwt.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
