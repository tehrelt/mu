package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

func Logout(auther authpb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, ok := c.Locals(middlewares.TokenLocalKey).(string)
		if !ok {
			slog.Error("failed to find token in fiber locals")
			return fiber.NewError(fiber.StatusInternalServerError, "unexpected error")
		}

		_, err := auther.Logout(c.UserContext(), &authpb.LogoutRequest{
			AccessToken: token,
		})
		if err != nil {
			slog.Error("failed to logout", "error", err)
			return fiber.NewError(fiber.StatusInternalServerError, "unexpected error")
		}

		return c.SendStatus(200)
	}
}
