package main

import (
	"context"
	"flag"

	"github.com/joho/godotenv"
	"github.com/tehrelt/mu/telegram-bot/internal/app"
)

var (
	envPath string
)

func init() {
	flag.StringVar(&envPath, "env", "", "path to .env file")
}

func main() {
	flag.Parse()

	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			panic(err)
		}
	}

	ctx := context.Background()
	app, fn, err := app.NewApp(ctx)
	if err != nil {
		panic(err)
	}
	defer fn()

	app.Run(ctx)
}
