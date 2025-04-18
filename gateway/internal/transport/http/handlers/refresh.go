package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type RefreshRequest struct {
	AccessToken string `json:"accessToken"`
}

func Refresh(auther authpb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("refresh_token", "")
		if token == "" {
			return fiber.NewError(401, "no refresh token")
		}

		res, err := auther.Refresh(c.UserContext(), &authpb.RefreshRequest{
			RefreshToken: token,
		})
		if err != nil {
			return err
		}

		ret := &RefreshRequest{
			AccessToken: res.Tokens.AccessToken,
		}

		c.Cookie(createCookie(res.Tokens.RefreshToken))

		return c.JSON(ret)
	}
}
