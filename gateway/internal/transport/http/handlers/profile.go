package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

func Profile(auther authpb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals(middlewares.TokenLocalKey).(string)
		if !ok {
			slog.Error("failed to get token from context")
			return fiber.ErrUnauthorized
		}
		resp, err := auther.Profile(c.UserContext(), &authpb.ProfileRequest{
			AccessToken: token,
		})
		if err != nil {
			return err
		}
		return c.JSON(resp)
	}
}
