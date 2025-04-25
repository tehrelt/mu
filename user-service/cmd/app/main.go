package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/tehrelt/mu/user-service/internal/app"
)

var (
	env string
)

func init() {
	flag.StringVar(&env, "env", "", "path to .env file")
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
		panic(err)
	}
	defer cleanup()

	instance.Run(ctx)
}
