package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	// "github.com/zalynskaya/murmur/internal/postgresql"
	// v1 "messenger-rest-api/app/internal/controllers/http/v1"

	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"

	"github.com/zalynskaya/murmur/internal/config"
)

type App struct {
	cfg *config.Config

	router     *httprouter.Router
	httpServer *http.Server
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	log.Println("router initializing")
	router := httprouter.New()

	log.Println("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	/*	pgConfig := postgresql.NewPgConfig(
			config.PostgreSQL.Username,
			config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
		)

		pgClient, err := postgresql.NewClient(ctx, 5, time.Second*5, pgConfig)
		if err != nil {
			log.Fatalln(err)
		}

		userStorage := postgresql2.NewUserStorage(pgClient)
		userService := user.NewUserService(userStorage)
		userHandler := v1.NewUserHandler(userService)
		userHandler.Register(router)

		chatStorage := postgresql2.NewChatStorage(pgClient)
		chatService := chat.NewChatService(chatStorage, userStorage)
		chatHandler := v1.NewChatHandler(chatService)
		chatHandler.Register(router)

		messageStorage := postgresql2.NewMessageStorage(pgClient)
		messageService := message.NewMessageService(messageStorage, userStorage, chatStorage)
		messageHandler := v1.NewMessageHandler(messageService)
		messageHandler.Register(router) */

	return App{
		cfg:    config,
		router: router,
		//pgClient: pgClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, _ := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP()
	})

	return grp.Wait()
}

func (a *App) startHTTP() error {
	log.Println("HTTP Server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		log.Fatalln("failed to create listener")
	}

	a.httpServer = &http.Server{
		Handler:      a.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Println("server shutdown")
		default:
			log.Fatal(err)
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return err
}
