package main

import (
	"context"
	"log"

	"github.com/zalynskaya/murmur/cmd/api_server/server"
	"github.com/zalynskaya/murmur/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("config initializing")
	cfg := config.GetConfig()

	a, err := server.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, err)
	}

	log.Println("Running Application")
	a.Run(ctx)
}
