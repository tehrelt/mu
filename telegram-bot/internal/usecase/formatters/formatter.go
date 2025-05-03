package formatters

import (
	"encoding/json"

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
	switch event.Header().EventType {
	case events.EventTicketStatusChanged:
		return f.ticketStatusFormatter(event.(*events.TicketStatusChanged))
	default:
		j, _ := json.Marshal(event)
		return string(j)
	}
}
