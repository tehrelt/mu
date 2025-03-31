package jwt

import "github.com/tehrelt/mu/auth-service/internal/config"

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type TokenConfig struct {
	Secret []byte
	TTL    int
}

type JwtClient struct {
	accessConfig  TokenConfig
	refreshConfig TokenConfig
}

func New(cfg *config.Config) *JwtClient {
	return &JwtClient{
		accessConfig: TokenConfig{
			Secret: []byte(cfg.Jwt.AccessSecret),
			TTL:    cfg.Jwt.AccessTTL,
		},
		refreshConfig: TokenConfig{
			Secret: []byte(cfg.Jwt.RefreshSecret),
			TTL:    cfg.Jwt.RefreshTTL,
		},
	}
}

func (jc JwtClient) GetTokenConfig(tokenType TokenType) TokenConfig {
	if tokenType == AccessToken {
		return jc.accessConfig
	}

	return jc.refreshConfig
}
