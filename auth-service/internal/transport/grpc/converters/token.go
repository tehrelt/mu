package converters

import (
	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/authpb"
)

func TokenPairToPb(pair *dto.TokenPair) *authpb.Tokens {
	return &authpb.Tokens{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}
}
