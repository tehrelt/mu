package formatters

import (
	"fmt"

	"github.com/tehrelt/mu/telegram-bot/internal/events"
)

func _status(raw string) string {
	switch raw {
	case "rejected":
		return "отклонен"
	case "approved":
		return "принят"
	case "pending":
		return "в ожидании"
	default:
		return "неизвестный"
	}
}

func (f *Formatter) ticketStatusFormatter(event *events.TicketStatusChanged) string {

	link := fmt.Sprintf("%s/tickets/%s", f.cfg.FrontendLink, event.TicketId)

	return fmt.Sprintf(
		`*ИЗМЕНЕНИЕ СТАТУСА ЗАЯВКИ*

Вашей [заявке](%s) присвоен новый статус: *%s*

Ссылка: %s`,
		link,
		_status(event.NewStatus),
		link,
	)
}
