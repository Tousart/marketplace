package main

import (
	"log"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/tousart/marketplace/config"
	"github.com/tousart/marketplace/middleware"
	"github.com/tousart/marketplace/repository/postgres"
	"github.com/tousart/marketplace/usecase/service"

	_ "github.com/lib/pq"
	handlers "github.com/tousart/marketplace/API/http"
	pfkHTTP "github.com/tousart/marketplace/pkg/http"
)

func main() {
	cfgPath := config.ParseFlag()
	cfg, err := config.MustLoad(cfgPath)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	// repository
	usersRepo, err := postgres.NewUsersRepo(cfg)
	if err != nil {
		log.Fatalf("failed to create users repository: %v", err)
	}

	advertsRepo, err := postgres.NewAdvertsRepo(cfg)
	if err != nil {
		log.Fatalf("failed to create adverts repository: %v", err)
	}

	// usecase
	mu := new(sync.Mutex)
	authService := service.NewAuthService(usersRepo)
	advertService := service.NewAdvertsService(advertsRepo, mu)

	// handlers
	userHandlers := handlers.NewUsersHandler(authService)
	advertHandlers := handlers.NewAdvertsHandler(advertService)

	// router
	r := chi.NewRouter()
	r.Use(middleware.Authorization)
	userHandlers.WithUsersHandlers(r)
	advertHandlers.WithAdvertsHandlers(r)

	log.Printf("server started on %s\n", cfg.HTTP.Address)
	if err := pfkHTTP.CreateAndRunServer(r, cfg); err != nil {
		log.Fatalf("failed to create and run server: %v", err)
	}
}
