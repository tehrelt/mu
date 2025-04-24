package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/register-service/internal/config"
	"github.com/tehrelt/mu/register-service/internal/services/registerservice"
)

var maxInterval = 60 * time.Minute

type Cron struct {
	cfg        *config.Config
	regservice *registerservice.Service
	interval   time.Duration
}

func (c *Cron) incrementInterval() {
	if c.interval < maxInterval {
		c.interval *= 2
	} else {
		c.interval = maxInterval
	}
}

func New(cfg *config.Config, regservice *registerservice.Service) *Cron {
	return &Cron{
		cfg:        cfg,
		regservice: regservice,
		interval:   5 * time.Second,
	}
}

func (c *Cron) Start(ctx context.Context) {

	t := time.NewTicker(c.interval)
	stopCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-t.C:
				if err := c.reg(ctx); err != nil {
					slog.Error("failed to create admin user", sl.Err(err))
					c.incrementInterval()
				}
			}
		}
	}()

	go func() {
		defer t.Stop()
		<-stopCh
	}()
}
