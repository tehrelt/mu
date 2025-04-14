package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/pkg/pb/registerpb"
	"google.golang.org/grpc/status"
)

type RegisterRequest struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	Phone      string `json:"phone"`
	Passport   struct {
		Number int32 `json:"number"`
		Series int32 `json:"series"`
	} `json:"passport"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Snils    string   `json:"snils"`
	Roles    []string `json:"roles,omitempty"`
}

type RegisterResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func Register(register registerpb.RegisterServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		roles := make([]registerpb.Role, 0)
		for _, role := range req.Roles {
			roles = append(roles, toProtoRole(role))
		}

		regResponse, err := register.Register(c.UserContext(), &registerpb.RegisterRequest{
			User: &registerpb.User{
				LastName:       req.LastName,
				FirstName:      req.FirstName,
				MiddleName:     req.MiddleName,
				Email:          req.Email,
				Phone:          req.Phone,
				PassportNumber: req.Passport.Number,
				PassportSeries: req.Passport.Series,
				Snils:          req.Snils,
				Password:       req.Password,
				Roles:          roles,
			},
		})
		if err != nil {
			if e, ok := status.FromError(err); ok {
				return fiber.NewError(fiber.StatusConflict, e.Message())
			}
			return err
		}

		res := &RegisterResponse{
			AccessToken:  regResponse.Tokens.AccessToken,
			RefreshToken: regResponse.Tokens.RefreshToken,
		}

		return c.JSON(res)
	}
}

func toProtoRole(role string) registerpb.Role {
	switch role {
	case "admin":
		return registerpb.Role_ROLE_ADMIN
	case "regular":
		return registerpb.Role_ROLE_REGULAR
	default:
		return registerpb.Role_ROLE_UNKNOWN
	}
}
