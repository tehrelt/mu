package formatters

import (
	"fmt"

	"github.com/tehrelt/mu/mailer/internal/events"
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
		`<p><strong>ИЗМЕНЕНИЕ СТАТУСА ЗАЯВКИ</strong></p>
<br/>
<p>Вашей <a href="%s">заявке</a> присвоен новый статус: <strong>%s</strong></p>
<br/>
С уважением, <a href="%s">команда Мои услуги</a>`,
		link,
		_status(event.NewStatus),
		f.cfg.FrontendLink,
	)
}
