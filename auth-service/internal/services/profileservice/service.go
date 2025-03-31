package profileservice

import (
	"context"

	"github.com/google/uuid"
	"github.com/tehrelt/mu/auth-service/internal/config"
	"github.com/tehrelt/mu/auth-service/internal/lib/jwt"
	"github.com/tehrelt/mu/auth-service/internal/models"
)

type UserProvider interface {
	UserById(ctx context.Context, id uuid.UUID) (*models.User, error)
}

type RoleProvider interface {
	Roles(ctx context.Context, userId uuid.UUID) ([]models.Role, error)
}

type ProfileService struct {
	cfg          *config.Config
	userProvider UserProvider
	roleProvider RoleProvider
	jc           *jwt.JwtClient
}

func New(
	cfg *config.Config,
	up UserProvider,
	jc *jwt.JwtClient,
	rp RoleProvider,
) *ProfileService {
	return &ProfileService{
		cfg:          cfg,
		userProvider: up,
		roleProvider: rp,
		jc:           jc,
	}
}
