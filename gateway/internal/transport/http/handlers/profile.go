package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/dto"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type ProfileResponse struct {
	Id         string `json:"id"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	Email      string `json:"email"`
}

func Profile(auther authpb.AuthServiceClient, accounter accountpb.AccountServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile, ok := c.Locals(middlewares.ProfileLocalKey).(*dto.UserProfile)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		resp := &ProfileResponse{
			Id:         profile.Id.String(),
			LastName:   profile.LastName,
			FirstName:  profile.FirstName,
			MiddleName: profile.MiddleName,
			Email:      profile.Email,
		}

		return c.JSON(resp)
	}
}
