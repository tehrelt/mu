package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func Login(auther authpb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		resp, err := auther.Login(c.UserContext(), &authpb.LoginRequest{
			Login: &authpb.LoginRequest_Email{
				Email: req.Login,
			},
			Password: req.Password,
		})
		if err != nil {
			return err
		}

		res := &LoginResponse{
			AccessToken:  resp.Tokens.AccessToken,
			RefreshToken: resp.Tokens.RefreshToken,
		}

		return c.JSON(res)
	}
}
