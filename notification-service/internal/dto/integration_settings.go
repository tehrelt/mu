package dto

type IntegrationSettings struct {
	UserId         string  `json:"id"`
	Email          string  `json:"email"`
	TelegramChatId *string `json:"telegramId"`
}

func NewIntegrationSettings(userId string, email string) *IntegrationSettings {
	return &IntegrationSettings{
		UserId: userId,
		Email:  email,
	}
}

func (setting *IntegrationSettings) WithTelegramChatId(chatId string) *IntegrationSettings {
	setting.TelegramChatId = new(string)
	*setting.TelegramChatId = chatId
	return setting
}
