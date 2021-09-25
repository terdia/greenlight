package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/terdia/greenlight/config"
	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/infrastructures/persistence/postgres"
	"github.com/terdia/greenlight/infrastructures/persistence/postgres/repository"
	"github.com/terdia/greenlight/internal/commons"
	"github.com/terdia/greenlight/internal/registry"
	"github.com/terdia/greenlight/src/movies/services"
)

type application struct {
	config   *config.Config
	registry registry.Registry
}

func main() {
	cfg := config.Get()
	logger := logger.NewLogger().Logger

	db, err := postgres.OpenDb(cfg.Db)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Printf("database connection pool established")

	services := registry.NewServices(
		services.NewMovieService(repository.NewMovieRepoitory(db)),
		commons.NewUtil(logger),
	)

	app := &application{
		config:   cfg,
		registry: registry.NewRegistry(db, services),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.AppPort),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.Env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
