package jwt

import (
	"github.com/tehrelt/moi-uslugi/auth-service/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	entity.UserClaims
	jwt.RegisteredClaims
}
