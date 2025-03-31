package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tehrelt/mu/auth-service/internal/dto"
)

func (jc *JwtClient) Verify(tokenString string, tokenType TokenType) (*dto.UserClaims, error) {

	cfg := jc.GetTokenConfig(tokenType)

	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, err
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return nil, errors.New("unable to parse claims")
	}

	return &claims.UserClaims, nil
}
