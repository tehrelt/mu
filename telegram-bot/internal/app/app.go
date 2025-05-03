package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/tehrelt/mu/telegram-bot/internal/config"
	"github.com/tehrelt/mu/telegram-bot/internal/transport/amqp"
	"github.com/tehrelt/mu/telegram-bot/internal/transport/tg"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	cfg *config.Config

	bot      *tg.Bot
	consumer *amqp.Consumer

	tracer trace.Tracer
}

func newApp(cfg *config.Config, bot *tg.Bot, c *amqp.Consumer, tracer trace.Tracer) *App {
	return &App{
		cfg:      cfg,
		bot:      bot,
		consumer: c,
		tracer:   tracer,
	}
}

func (a *App) Run(ctx context.Context) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	wg := sync.WaitGroup{}

	now := time.Now()

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.bot.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.consumer.Run(ctx); err != nil {
			slog.Error("consumer error", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()
	slog.Info("application stopping...", slog.Duration("uptime", time.Since(now)))
	wg.Wait()
}
