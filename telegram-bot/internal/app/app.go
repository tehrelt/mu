package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/tehrelt/mu/telegram-bot/internal/config"
	"github.com/tehrelt/mu/telegram-bot/internal/transport/tg"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	cfg    *config.Config
	bot    *tg.Bot
	tracer trace.Tracer
}

func newApp(cfg *config.Config, bot *tg.Bot, tracer trace.Tracer) *App {
	return &App{
		cfg:    cfg,
		bot:    bot,
		tracer: tracer,
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

	<-ctx.Done()
	slog.Info("application stopping...", slog.Duration("uptime", time.Since(now)))
	wg.Wait()
}
