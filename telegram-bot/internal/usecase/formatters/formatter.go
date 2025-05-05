package formatters

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/tehrelt/mu/telegram-bot/internal/config"
	"github.com/tehrelt/mu/telegram-bot/internal/events"
)

type Formatter struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Formatter {
	return &Formatter{
		cfg: cfg,
	}
}

func (f *Formatter) Format(event events.Event) string {
	return f.escape(f.format(event))
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
		"-", "\\-",
		">", "\\>",
		"<", "\\<",
		".", "\\.",
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
		`*%s* %s %.2fр.
%s: %.2fр.`,
		event.Address, header, float64(delta)/100,
		balance, float64(event.NewBalance)/100,
	)
}
