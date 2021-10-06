package main

import (
	"os"

	_ "github.com/lib/pq"

	"github.com/terdia/greenlight/config"
	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/infrastructures/persistence/postgres"
	"github.com/terdia/greenlight/internal/registry"
)

type application struct {
	config   *config.Config
	registry registry.Registry
	logger   *logger.Logger
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
		logger:   logger,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
