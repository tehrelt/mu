package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/tehrelt/mu/auth-service/internal/app"

	"github.com/joho/godotenv"
)

var (
	env string
)

func init() {
	flag.StringVar(&env, "env", "", "path to env file")
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
