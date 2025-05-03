package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/notificationpb"
)

type TelegramOTPResponse struct {
	OTP    string `json:"otp"`
	UserId string `json:"userId"`
}

func GetTelegramOTP(notifier notificationpb.NotificationServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		user, _ := middlewares.UserFromLocals(c)

		otpresp, err := notifier.TelegramOtp(ctx, &notificationpb.TelegramOtpRequest{
			UserId: user.Id.String(),
		})
		if err != nil {
			return err
		}

		res := TelegramOTPResponse{
			OTP:    otpresp.Otp,
			UserId: user.Id.String(),
		}

		return c.JSON(res)
	}
}
