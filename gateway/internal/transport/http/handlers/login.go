package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Login(auther authpb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		slog.Debug("login request", slog.Any("payload", req))
		resp, err := auther.Login(c.Context(), &authpb.LoginRequest{
			Login: &authpb.LoginRequest_Email{
				Email: req.Login,
			},
			Password: req.Password,
		})
		if err != nil {
			return err
		}

		return c.JSON(&fiber.Map{
			"accessToken":  resp.Tokens.AccessToken,
			"refreshToken": resp.Tokens.RefreshToken,
		})
	}
}
