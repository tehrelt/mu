package usersapi

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/storage"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/userspb"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (api *Api) UserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	log := slog.With(sl.Method("usersapi.UserById"), slog.String("userId", userId.String()))

	log.Debug("searching for user")
	resp, err := api.client.Find(ctx, &userspb.FindRequest{
		SearchBy: &userspb.FindRequest_Id{Id: userId.String()},
	})
	if err != nil {
		log.Warn("failed to search for user", sl.Err(err))

		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				log.Debug("user not found")
				return nil, storage.ErrUserNotFound
			}
		}

		return nil, err
	}
	user := userFromProto(resp.User)

	return user, nil
}

func (api *Api) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	log := slog.With(sl.Method("usersapi.UserByEmail"), slog.String("email", email))

	log.Debug("searching for user")
	resp, err := api.client.Find(ctx, &userspb.FindRequest{
		SearchBy: &userspb.FindRequest_Email{Email: email},
	})
	if err != nil {
		log.Warn("failed to search for user", sl.Err(err))
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				log.Debug("user not found")
				return nil, storage.ErrUserNotFound
			}
		}
		return nil, err
	}
	user := userFromProto(resp.User)

	return user, nil
}
