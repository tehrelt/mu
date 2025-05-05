package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tehrelt/mu/gateway/internal/transport/http/middlewares"
	"github.com/tehrelt/mu/gateway/pkg/pb/notificationpb"
)

type IntegrationsSettings struct {
	HasTelegram bool `json:"hasTelegram"`
}

func GetIntegrationsSettings(notifier notificationpb.NotificationServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		user, _ := middlewares.UserFromLocals(c)

		settings, err := notifier.Integrations(ctx, &notificationpb.IntegrationsRequest{
			UserId: user.Id.String(),
		})
		if err != nil {
			return err
		}

		resp := IntegrationsSettings{
			HasTelegram: settings.GetTelegramChatId() != "",
		}

		return c.JSON(resp)
	}
}

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
