package jwt

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tehrelt/mu/auth-service/internal/dto"
)

func (jc *JwtClient) Sign(user *dto.UserClaims, tokenType TokenType) (string, error) {
	cfg := jc.GetTokenConfig(tokenType)

	ttl := time.Duration(cfg.TTL) * time.Minute
	iat := jwt.NewNumericDate(time.Now())
	exp := jwt.NewNumericDate(time.Now().Add(ttl))
	slog.Info(
		"signing token",
		slog.String("token_type", tokenType.String()),
		slog.Duration("ttl", ttl),
		slog.Any("exp", exp),
		slog.Any("iat", iat),
	)
	payload := claims{
		UserClaims: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: exp,
			IssuedAt:  iat,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signed, err := token.SignedString(cfg.Secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}
