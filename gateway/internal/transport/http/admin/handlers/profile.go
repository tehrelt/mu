package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/tehrelt/mu/gateway/internal/dto"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/accountpb"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
)

type ProfileRequest struct {
}

type ProfileResponse struct {
	Id         string     `json:"id"`
	LastName   string     `json:"lastName"`
	FirstName  string     `json:"firstName"`
	MiddleName string     `json:"middleName"`
	Email      string     `json:"email"`
	Roles      []dto.Role `json:"roles"`
}

func Profile(auther authpb.AuthServiceClient, accounter accountpb.AccountServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		profile, ok := c.Locals(middlewares.ProfileLocalKey).(*dto.UserProfile)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		expectedRole := dto.Role(c.Query("role", ""))
		if expectedRole != "" {
			slog.Debug(
				"role param is set",
				slog.String("expectedRole", string(expectedRole)),
				slog.Bool("isValid", expectedRole.Validate()),
				slog.Any("accRoles", profile.Roles),
			)
			if expectedRole.Validate() {
				if !lo.Contains(profile.Roles, expectedRole) {
					return fiber.NewError(fiber.StatusForbidden, "forbidden")
				}
			}
		}

		resp := &ProfileResponse{
			Id:         profile.Id.String(),
			LastName:   profile.LastName,
			FirstName:  profile.FirstName,
			MiddleName: profile.MiddleName,
			Email:      profile.Email,
			Roles:      profile.Roles,
		}

		return c.JSON(resp)
	}
}
