package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type RefreshRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func Refresh(auther authpb.AuthServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, ok := c.Locals(middlewares.TokenLocalKey).(string)
		if !ok {
			return fiber.NewError(401, "where token")
		}

		res, err := auther.Refresh(c.UserContext(), &authpb.RefreshRequest{
			RefreshToken: token,
		})
		if err != nil {
			return err
		}

		ret := &RefreshRequest{
			AccessToken:  res.Tokens.AccessToken,
			RefreshToken: res.Tokens.RefreshToken,
		}

		return c.JSON(ret)
	}
}
