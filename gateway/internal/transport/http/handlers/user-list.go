package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu-lib/tracer"
	"github.com/tehrelt/mu/gateway/pkg/pb/userpb"
	"go.opentelemetry.io/otel"
)

type UserSnippet struct {
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
	Users []UserSnippet `json:"users"`
	Total uint64        `json:"total"`
}

func UserListHandler(userapi userpb.UserServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.UserContext()
		var req UserListRequest

		if err := c.QueryParser(&req); err != nil {
			return err
		}

		slog.Debug("user list request", slog.Any("filters", req))

		ctx, span := otel.Tracer(tracer.TracerKey).Start(ctx, fmt.Sprintf("user list request"))
		defer span.End()

		stream, err := userapi.List(ctx, &userpb.ListRequest{
			Offset: (req.Page - 1) * req.Limit,
			Limit:  req.Limit,
		})
		if err != nil {
			return err
		}

		var users []UserSnippet

		for {
			resp, err := stream.Recv()
			span.AddEvent("response received")
			if err == io.EOF {
				span.AddEvent("stream ended")
				break
			}
			if err != nil {
				return err
			}

			slog.Debug("recieved users chunk", slog.Int("chunk size", len(resp.UsersChunk)))

			for _, user := range resp.UsersChunk {
				u := UserSnippet{
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

			span.AddEvent("chunkprocessed")
		}

		return c.JSON(&ListUsersResponse{
			Users: users,
		})
	}
}
