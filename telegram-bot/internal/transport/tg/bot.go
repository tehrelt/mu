package tg

import (
	"context"
	"log/slog"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/telegram-bot/internal/config"
	"gopkg.in/telebot.v4"
)

type Bot struct {
	cfg    *config.Config
	bot    *telebot.Bot
	logger *slog.Logger
}

func New(cfg *config.Config, bot *telebot.Bot) *Bot {
	instance := &Bot{
		cfg:    cfg,
		bot:    bot,
		logger: slog.With(sl.Module("TgBot")),
	}

	instance.setup()

	return instance
}

func (b *Bot) setup() {
	b.bot.Handle("/start", b.startHandler())
}

func (b *Bot) Run(ctx context.Context) {
	b.logger.Info("starting up bot")
	go b.bot.Start()

	<-ctx.Done()
	b.bot.Stop()
}
