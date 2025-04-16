package middlewares

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func Cookies() fiber.Handler {
	return func(c *fiber.Ctx) error {
		slog.Info("cookies of request", slog.Any("cookies", c.Cookies("accessToken", "empty accessToken")))
		return c.Next()
	}
}
