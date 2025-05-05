package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func createCookie(token string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HTTPOnly: true,
		SameSite: "lax",
	}
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	// RefreshToken string `json:"refreshToken"`
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
			if e, ok := status.FromError(err); ok {
				if e.Code() == codes.Unauthenticated {
					return fiber.NewError(fiber.StatusUnauthorized, e.Message())
				}
			}
			return err
		}

		res := &LoginResponse{
			AccessToken: resp.Tokens.AccessToken,
		}

		c.Cookie(createCookie(resp.Tokens.RefreshToken))

		return c.JSON(res)
	}
}
