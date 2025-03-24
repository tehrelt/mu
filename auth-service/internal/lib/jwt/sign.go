package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
)

func (jc *JwtClient) Sign(user *dto.UserClaims, tokenType TokenType) (string, error) {
	cfg := jc.GetTokenConfig(tokenType)

	payload := claims{
		UserClaims: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.TTL))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signed, err := token.SignedString(cfg.Secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}
