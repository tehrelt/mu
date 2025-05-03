package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/tehrelt/mu-lib/sl"
	"github.com/tehrelt/mu/notification-service/internal/app"
)

var (
	env string
)

func init() {
	flag.StringVar(&env, "env", "", "environment")
}

func main() {
	flag.Parse()

	if env != "" {
		if err := godotenv.Load(env); err != nil {
			panic(fmt.Errorf("cannot load env: %w", err))
		}
	}

	ctx := context.Background()

	instance, cleanup, err := app.New(ctx)
	if err != nil {
		slog.Error("failed to start application", sl.Err(err))
		panic(err)
	}
	defer cleanup()

	instance.Run(ctx)
}
