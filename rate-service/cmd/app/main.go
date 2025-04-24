package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/tehrelt/mu/rate-service/internal/app"
)

var (
	local bool
)

func init() {
	flag.BoolVar(&local, "local", false, "run in local mode")
}

func main() {
	flag.Parse()

	if local {
		if err := godotenv.Load(); err != nil {
			panic(fmt.Errorf("cannot load env: %w", err))
		}
	}

	ctx := context.Background()

	instance, cleanup, err := app.New()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	instance.Run(ctx)
}
