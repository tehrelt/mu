package models

import "github.com/google/uuid"

type Integration struct {
	UserId         uuid.UUID
	TelegramChatId *string
}

func (i *Integration) SetTelegramChatId(chatId string) {
	i.TelegramChatId = new(string)
	*i.TelegramChatId = chatId
}
