package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

func Refresh(auther authpb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(500).SendString("unimplemented")
	}
}
