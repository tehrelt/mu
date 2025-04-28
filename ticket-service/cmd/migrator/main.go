package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/tehrelt/mu/ticket-service/internal/config"
	"github.com/tehrelt/mu/ticket-service/internal/storage/mongo"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			fmt.Println("Error loading .env file:", err)
			os.Exit(1)
		}
	}

	ctx := context.Background()

	cfg := config.New()

	client, err := mongodriver.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.ConnectionString()))
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		return
	}

	db := client.Database(cfg.Mongo.Database)
	if err := db.CreateCollection(ctx, mongo.TICKETS_COLLECTION); err != nil {
		fmt.Println("Error creating collection:", err)
		return
	}
}
