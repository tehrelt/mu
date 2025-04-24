package handlers

import (
	"io"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/userpb"
)

type User struct {
	Id         string     `json:"id"`
	LastName   string     `json:"lastName"`
	FirstName  string     `json:"firstName"`
	MiddleName string     `json:"middleName"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type UserListRequest struct {
	Page  uint64 `query:"page"`
	Limit uint64 `query:"limit"`
}

type ListUsersResponse struct {
	Users []User `json:"users"`
}

func UserListHandler(userapi userpb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var req UserListRequest

		if err := c.QueryParser(&req); err != nil {
			return err
		}

		slog.Debug("user list request", slog.Any("filters", req))

		stream, err := userapi.List(c.Context(), &userpb.ListRequest{
			Offset: (req.Page - 1) * req.Limit,
			Limit:  req.Limit,
		})
		if err != nil {
			return err
		}

		var users []User

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			slog.Debug("recieved users chunk", slog.Any("chunk", resp.UsersChunk))

			for _, user := range resp.UsersChunk {
				u := User{
					Id:         user.Id,
					LastName:   user.Fio.Lastname,
					FirstName:  user.Fio.Firstname,
					MiddleName: user.Fio.Middlename,
					Email:      user.Email,
					Phone:      user.PersonalData.Phone,
					CreatedAt:  time.Unix(user.CreatedAt, 0),
				}

				if user.UpdatedAt != 0 {
					*u.UpdatedAt = time.Unix(user.UpdatedAt, 0)
				}

				users = append(users, u)
			}
		}

		return c.JSON(&ListUsersResponse{
			Users: users,
		})
	}
}
