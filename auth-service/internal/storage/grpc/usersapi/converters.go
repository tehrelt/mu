package usersapi

import (
	"github.com/tehrelt/moi-uslugi/auth-service/internal/models"
	"github.com/tehrelt/moi-uslugi/auth-service/pkg/pb/userspb"
)

func userFromProto(user *userspb.User) *models.User {
	return &models.User{
		Id:    user.Id,
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
	}
}
