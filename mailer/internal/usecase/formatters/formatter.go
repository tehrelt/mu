package formatters

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/tehrelt/mu/mailer/internal/config"
	"github.com/tehrelt/mu/mailer/internal/events"
)

type Formatter struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Formatter {
	return &Formatter{
		cfg: cfg,
	}
}

func (f *Formatter) Format(event events.Event) (header string, body string) {
	return f.formatHeader(event), f.escape(f.format(event))
}

func (f *Formatter) formatHeader(event events.Event) string {

	switch e := event.(type) {
	case *events.TicketStatusChanged:
		return fmt.Sprintf("[Мои услуги] новый статус заявки %s", _status(e.NewStatus))
	case *events.BalanceChanged:

		header := "ПОСТУПЛЕНИЕ"
		delta := e.Delta()

		if delta < 0 {
			header = "СПИСАНИЕ"
			delta = -delta
		}

		return fmt.Sprintf("[Мои услуги] %s по адресу %s", header, e.Address)
	default:
		return fmt.Sprintf("[Мои услуги] уведомление от %s", event.Header().Timestamp.String())
	}
}

func (f *Formatter) format(event events.Event) string {

	slog.Debug("formating", slog.Any("event", event))

	switch event.Header().EventType {
	case events.EventTicketStatusChanged:
		return f.ticketStatusFormatter(event.(*events.TicketStatusChanged))
	case events.EventBalanceChanged:
		return f.balanceChangedFormatter(event.(*events.BalanceChanged))
	default:
		j, _ := json.Marshal(event)
		return string(j)
	}
}

func (f *Formatter) escape(s string) string {
	replacer := strings.NewReplacer(
	// "-", "\\-",
	// ">", "\\>",
	// "<", "\\<",
	// ".", "\\.",
	)
	return replacer.Replace(s)
}

func (f *Formatter) balanceChangedFormatter(event *events.BalanceChanged) string {

	header := "ПОСТУПЛЕНИЕ"
	delta := event.Delta()
	balance := "Доступно"

	if delta < 0 {
		header = "СПИСАНИЕ"
		delta = -delta
	}

	if event.NewBalance < 0 {
		balance = "Задолженность"
	}

	return fmt.Sprintf(
		`<p>%s %.2fр.</p>
		<p>Адрес: <strong>%s</strong></p>
<p>%s: %.2fр.</p>`,
		header, float64(delta)/100,
		event.Address,
		balance, float64(event.NewBalance)/100,
	)
}
