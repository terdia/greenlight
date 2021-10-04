package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/terdia/greenlight/config"
	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/infrastructures/persistence/postgres"
	"github.com/terdia/greenlight/internal/registry"
)

type application struct {
	config   *config.Config
	registry registry.Registry
}

//https://greenlight.docker.local/v1
// @title Greenlight API documentation
// @version 1.0.0
// @host localhost:4000
// @BasePath /v1
func main() {
	cfg := config.NewConfig()
	logger := logger.New(os.Stdout, logger.LevelInfo)

	db, err := postgres.OpenDb(cfg.Db)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config:   cfg,
		registry: registry.NewRegistry(db, logger),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.AppPort),
		Handler:      app.routes(),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.PrintInfo("starting server", map[string]string{
		"addr": fmt.Sprint(cfg.AppPort),
		"env":  cfg.Env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}
