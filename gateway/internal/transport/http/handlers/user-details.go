package handlers

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	Id             string     `json:"id"`
	LastName       string     `json:"lastName"`
	FirstName      string     `json:"firstName"`
	MiddleName     string     `json:"middleName"`
	Email          string     `json:"email"`
	Phone          string     `json:"phone"`
	PassportSeries int32      `json:"passportSeries"`
	PassportNumber int32      `json:"passportNumber"`
	Snils          string     `json:"snils"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
}

func UserDetailHandler(userapi userpb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		id := c.Params("id", "")
		if id == "" {
			return fiber.ErrBadRequest
		}

		slog.Debug("user find request", slog.Any("id", id))
		resp, err := userapi.Find(c.UserContext(), &userpb.FindRequest{
			SearchBy: &userpb.FindRequest_Id{
				Id: id,
			},
		})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				if e.Code() == codes.NotFound {
					return fiber.ErrNotFound
				}
			}
			return err
		}

		user := resp.User

		u := User{
			Id:             user.Id,
			LastName:       user.Fio.Lastname,
			FirstName:      user.Fio.Firstname,
			MiddleName:     user.Fio.Middlename,
			Email:          user.Email,
			Phone:          user.PersonalData.Phone,
			PassportSeries: user.PersonalData.Passport.Series,
			PassportNumber: user.PersonalData.Passport.Number,
			Snils:          user.PersonalData.Snils,
			CreatedAt:      time.Unix(user.CreatedAt, 0),
		}

		return c.JSON(u)
	}
}
