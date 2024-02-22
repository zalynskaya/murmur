package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/zalynskaya/murmur/internal/database"
	"github.com/zalynskaya/murmur/internal/repo"
	"github.com/zalynskaya/murmur/internal/service"
	"golang.org/x/sync/errgroup"

	_ "github.com/zalynskaya/murmur/cmd/api-server/docs"

	"github.com/zalynskaya/murmur/internal/config"
	v1 "github.com/zalynskaya/murmur/internal/v1"
)

type App struct {
	cfg *config.Config

	router     *httprouter.Router
	httpServer *http.Server
	pgClient   *pgxpool.Pool
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	log.Println("router initializing")
	router := httprouter.New()

	log.Println("swagger docs initializing")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	log.Println("database initializing")
	pgConfig := database.NewPgConfig(
		config.PostgreSQL.Username, config.PostgreSQL.Password,
		config.PostgreSQL.Host, config.PostgreSQL.Port, config.PostgreSQL.Database,
	)

	pgClient, err := database.NewClient(ctx, 5, time.Second*5, pgConfig)
	if err != nil {
		log.Fatalln(err)
	}

	userStorage := repo.NewUserStorage(pgClient)
	userService := service.NewUserService(userStorage)
	userHandler := v1.NewUserHandler(userService)
	userHandler.Register(router)

	return App{
		cfg:      config,
		router:   router,
		pgClient: pgClient,
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
