package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/tehrelt/mu/register-service/pkg/sl"
)

type Server interface {
	Run(ctx context.Context) error
}

type App struct {
	servers []Server
}

func newApp(servers []Server) *App {
	return &App{
		servers: servers,
	}
}

func (a *App) Register(s Server) {
	a.servers = append(a.servers, s)
}

func (a *App) Run(ctx context.Context) {

	wg := sync.WaitGroup{}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

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
