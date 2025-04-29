package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/tehrelt/mu/rate-service/internal/app"
)

var (
	envPath string
)

func init() {
	flag.StringVar(&envPath, "env", "", "path to env file")
}

func main() {
	flag.Parse()

	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			panic(fmt.Errorf("cannot load env: %w", err))
		}
	}

	ctx := context.Background()

	instance, cleanup, err := app.New(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	instance.Run(ctx)
}
