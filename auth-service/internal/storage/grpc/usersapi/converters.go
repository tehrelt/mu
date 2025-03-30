package usersapi

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/userpb"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/sl"
)

func userFromProto(user *userpb.User) (*models.User, error) {

	id, err := uuid.Parse(user.Id)
	if err != nil {
		slog.Error("failed to parse user id from proto", sl.Err(err))
		return nil, err
	}

	return &models.User{
		Id:    id,
		Email: user.Email,
		Fio: models.Fio{
			LastName:   user.Fio.Lastname,
			FirstName:  user.Fio.Firstname,
			MiddleName: user.Fio.Middlename,
		},
		PersonalData: models.PersonalData{
			Phone: user.PersonalData.Phone,
			Snils: user.PersonalData.Snils,
			Passport: models.Passport{
				Number: int(user.PersonalData.Passport.Number),
				Series: int(user.PersonalData.Passport.Series),
			},
		},
	}, nil
}
