package authservice

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/config"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/lib/jwt"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.0 --name=UserProvider
type UserProvider interface {
	UserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.0 --name=UserSaver
type UserCreator interface {
	Create(ctx context.Context, user *dto.CreateUser) (uuid.UUID, error)
}

type CredentialsSaver interface {
	Save(ctx context.Context, creds *models.Credentials) error
}

type CredentialsProvider interface {
	Password(ctx context.Context, userId uuid.UUID) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.0 --name=SessionsStorage
type SessionsStorage interface {
	Check(ctx context.Context, userId uuid.UUID, token string) error
	Save(ctx context.Context, userId uuid.UUID, token string) error
	Delete(ctx context.Context, userId uuid.UUID) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.0 --name=RoleStorage
type RoleStorage interface {
	Add(ctx context.Context, dto *dto.UserRoles) error
	Roles(ctx context.Context, userId uuid.UUID) ([]models.Role, error)
}

type AuthService struct {
	cfg                 *config.Config
	logger              *slog.Logger
	userSaver           UserCreator
	userProvider        UserProvider
	credentialSaver     CredentialsSaver
	credentialsProvider CredentialsProvider
	roles               RoleStorage
	sessions            SessionsStorage
	jwtClient           *jwt.JwtClient
}

func New(
	usaver UserCreator,
	uprovider UserProvider,
	r RoleStorage,
	s SessionsStorage,
	cfg *config.Config,
	jc *jwt.JwtClient,
	cs CredentialsSaver,
	cp CredentialsProvider,
) *AuthService {
	svc := &AuthService{
		cfg:                 cfg,
		logger:              slog.With(sl.Module("authservice.AuthService")),
		userSaver:           usaver,
		userProvider:        uprovider,
		roles:               r,
		sessions:            s,
		jwtClient:           jc,
		credentialSaver:     cs,
		credentialsProvider: cp,
	}

	return svc
}

// func (s *AuthService) setup(ctx context.Context) error {
// 	fn := "authservice.setup"
// 	log := s.logger.With(sl.Method(fn))

// 	log.Info("setup")

// 	email := s.cfg.DefaultAdmin.Email
// 	password := s.cfg.DefaultAdmin.Password

// 	_, err := s.userProvider.UserByEmail(ctx, email)
// 	if err == nil {
// 		return nil
// 	}

// 	if !errors.Is(err, storage.ErrUserNotFound) {
// 		return err
// 	}

// 	log.Info("root user not found")

// 	if _, err := s.Register(ctx, &dto.CreateUser{
// 		Email:    email,
// 		Password: password,
// 	}); err != nil {
// 		return err
// 	}

// 	user, err := s.userProvider.UserByEmail(ctx, email)
// 	if err != nil {
// 		return err
// 	}

// 	if err := s.roles.Add(
// 		ctx,
// 		&dto.AddRoles{
// 			UserId: user.Id,
// 			Roles:  []models.Role{models.Role_Admin},
// 		},
// 	); err != nil {
// 		return err
// 	}

// 	log.Info("root user created")

// 	return nil
// }
