package main

import (
	"context"
	"log"
	"os"

	"basic/app/config"
	"basic/di"
)

func main() {
	container, err := di.NewContainer(config.Params{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// get published service from di container
	server, err := container.Server(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	container.Close()
}
