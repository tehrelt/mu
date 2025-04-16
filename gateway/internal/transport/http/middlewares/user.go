package middlewares

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/gateway/internal/dto"
	"github.com/tehrelt/mu/gateway/pkg/pb/authpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RoleHandler func(roles ...dto.Role) fiber.Handler

func Auth(auther authpb.AuthServiceClient) RoleHandler {
	return func(roles ...dto.Role) fiber.Handler {
		return func(c *fiber.Ctx) error {

			token := c.Locals(TokenLocalKey).(string)

			if len(roles) != 0 {
				_, err := auther.Authorize(c.UserContext(), &authpb.AuthorizeRequest{
					Token: token,
					Roles: lo.Map(
						roles,
						func(role dto.Role, _ int) authpb.Role {
							return role.ToProto()
						},
					),
				})
				if err != nil {
					if e, ok := status.FromError(err); ok {
						if e.Code() == codes.PermissionDenied {
							slog.Warn("permission denied")
							return fiber.ErrForbidden
						}
						slog.Error("unexpected rpc error", sl.Err(err))
						return err
					}

					slog.Error("unexpected error", sl.Err(err))
					return err
				}
			}

			profileResponse, err := auther.Profile(c.UserContext(), &authpb.ProfileRequest{
				AccessToken: token,
			})
			if err != nil {
				if e, ok := status.FromError(err); ok {
					if e.Code() == codes.Unauthenticated {
						return fiber.NewError(fiber.ErrUnauthorized.Code, "unauthorized")
					}
					slog.Error("unexpected rpc error", sl.Err(err))
				} else {
					slog.Error("unexpected error", sl.Err(err))
				}
				return err
			}

			profile := &dto.UserProfile{}
			if err := profile.FromProto(profileResponse); err != nil {
				return err
			}

			c.Locals(ProfileLocalKey, profile)

			return c.Next()
		}
	}
}
