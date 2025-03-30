package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"

	"github.com/tehrelt/mu/user-service/pkg/sl"
	"go.opentelemetry.io/otel/trace"
)

type Server interface {
	Run(ctx context.Context) error
}

type App struct {
	servers []Server
	tracer  trace.Tracer
}

func newApp(servers []Server, t trace.Tracer) *App {
	return &App{
		servers: servers,
		tracer:  t,
	}
}

func (a *App) Register(s Server) {
	a.servers = append(a.servers, s)
}

func (a *App) Run(ctx context.Context) {

	wg := sync.WaitGroup{}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill)

	nctx, cancel := context.WithCancel(ctx)

	for _, server := range a.servers {
		wg.Add(1)
		go func(s Server) {
			defer wg.Done()
			err := s.Run(nctx)
			if err != nil {
				cancel()
				slog.Error("server run error", sl.Err(err))
			}
		}(server)
	}

	go func() {
		<-sigchan
		cancel()
	}()

	wg.Wait()

	<-nctx.Done()
	cancel()
	slog.Info("shutting down")

}
