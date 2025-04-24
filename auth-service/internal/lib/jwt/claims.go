package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/tehrelt/mu/auth-service/internal/dto"
)

type claims struct {
	dto.UserClaims
	jwt.RegisteredClaims
}
