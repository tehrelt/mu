package middlewares

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func BearerToken() fiber.Handler {
	return func(c *fiber.Ctx) error {

		header := c.Get("Authorization")

		if header == "" {
			return fiber.NewError(401, "no authorization header")
		}

		parts := strings.Split(header, " ")

		if strings.Compare(parts[0], "Bearer") != 0 {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		}

		slog.Debug("Bearer token", slog.String("token", parts[1]))

		c.Locals(TokenLocalKey, parts[1])

		return c.Next()
	}
}
