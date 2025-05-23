package usersapi

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/auth-service/internal/models"
	"github.com/tehrelt/mu/auth-service/internal/storage"
	"github.com/tehrelt/mu/auth-service/pkg/pb/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (api *Api) UserById(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	log := slog.With(sl.Method("usersapi.UserById"), slog.String("userId", userId.String()))

	log.Debug("searching for user")
	resp, err := api.client.Find(ctx, &userpb.FindRequest{
		SearchBy: &userpb.FindRequest_Id{Id: userId.String()},
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
	user, err := userFromProto(resp.User)
	if err != nil {
		log.Error("failed convert user from proto", sl.Err(err))
		return nil, err
	}

	return user, nil
}

func (api *Api) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	log := slog.With(sl.Method("usersapi.UserByEmail"), slog.String("email", email))

	log.Debug("searching for user")
	resp, err := api.client.Find(ctx, &userpb.FindRequest{
		SearchBy: &userpb.FindRequest_Email{Email: email},
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
	user, err := userFromProto(resp.User)
	if err != nil {
		log.Error("failed convert user from proto", sl.Err(err))
		return nil, err
	}

	return user, nil
}
