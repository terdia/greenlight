package main

import (
	"os"
	"sync"

	_ "github.com/lib/pq"

	"github.com/terdia/greenlight/config"
	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/infrastructures/persistence/postgres"
	"github.com/terdia/greenlight/internal/mailer"
	"github.com/terdia/greenlight/internal/registry"
)

type application struct {
	config   *config.Config
	registry registry.Registry
	logger   *logger.Logger
	wg       *sync.WaitGroup
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

	mailer := mailer.New(cfg.Smtp)

	wg := new(sync.WaitGroup)

	app := &application{
		config:   cfg,
		registry: registry.NewRegistry(db, logger, mailer, wg),
		logger:   logger,
		wg:       wg,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
