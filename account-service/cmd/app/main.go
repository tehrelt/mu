package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/tehrelt/mu/account-service/internal/app"
	"github.com/tehrelt/mu/account-service/pkg/sl"
)

var (
	env string
)

func init() {
	flag.StringVar(&env, "env", "dev", "environment")
}

func main() {
	flag.Parse()

	if env != "" {
		if err := godotenv.Load(env); err != nil {
			panic(fmt.Errorf("cannot load env: %w", err))
		}
	}

	ctx := context.Background()

	instance, cleanup, err := app.New()
	if err != nil {
		slog.Error("failed to start application", sl.Err(err))
		panic(err)
	}
	defer cleanup()

	instance.Run(ctx)
}
