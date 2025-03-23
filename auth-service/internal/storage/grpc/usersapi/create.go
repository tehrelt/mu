package usersapi

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/dto"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/userspb"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

func (api *Api) Create(ctx context.Context, req *dto.CreateUser) (uuid.UUID, error) {
	log := slog.With(sl.Method("userapi.Create"))

	log.Debug("sending create request", slog.Any("user request", req))
	resp, err := api.client.Create(ctx, &userspb.CreateRequest{
		Fio: &userspb.FIO{
			Lastname:   req.LastName,
			Firstname:  req.FirstName,
			Middlename: req.MiddleName,
		},
		Email: req.Email,
		PersonalData: &userspb.PersonalData{
			Passport: &userspb.Passport{
				Number: int32(req.Passport.Number),
				Series: int32(req.Passport.Series),
			},
			Snils: req.Snils,
			Phone: req.Phone,
		},
	})
	if err != nil {
		log.Warn("create user error", sl.Err(err))
		return uuid.Nil, err
	}

	log.Debug("user created", slog.Any("response", resp))
	id, err := uuid.Parse(resp.Id)
	if err != nil {
		log.Error("failed parse uuid", sl.Err(err))
		return uuid.Nil, err
	}

	return id, err
}
