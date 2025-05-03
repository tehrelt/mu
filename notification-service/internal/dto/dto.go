package dto

import (
	"github.com/google/uuid"
	"github.com/tehrelt/mu/notification-service/internal/models/otp"
)

type LinkTelegram struct {
	UserId         uuid.UUID
	Code           otp.OTP
	TelegramChatId string
}
