package main

import (
	"context"
	"log"

	"github.com/zalynskaya/murmur/cmd/api-server/server"
	"github.com/zalynskaya/murmur/internal/config"
)

// @title           Swagger API
// @version         1.0
// @description     This is a messenger service.

// @host      localhost:9000
// @BasePath  /

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
