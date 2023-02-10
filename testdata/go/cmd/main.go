package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/username/repo/http"
	"github.com/username/repo/postgres"
	"github.com/username/repo/redis"

	_ "github.com/lib/pq"
)

func main() {
	ctx := context.TODO()
	// Postgres
	pgConfig := postgres.Config{
		Password: os.Getenv("POSTGRES_PASSWORD"),
		User:     os.Getenv("POSTGRES_USER"),
		Name:     os.Getenv("POSTGRES_DB"),
		Host:     os.Getenv("POSTGRES_HOST"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	// Open psql
	psql, err := postgres.Open(pgConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer psql.Close()

	thingService := postgres.NewThingService(psql)

	// Redis
	redis := redis.Open()
	pong, err := redis.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pong)

	port := os.Getenv("GO_PORT")

	httpConfig := http.Config{
		Addr:         port,
		ThingService: thingService,
	}

	server := http.New(httpConfig)

	server.RegisterThingRoutes(ctx)

	env := os.Getenv("GO_ENV")

	if env == "development" {
		log.Fatal(server.Open())
	}
}
