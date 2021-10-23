package main

import (
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/terdia/greenlight/config"
	"github.com/terdia/greenlight/infrastructures/logger"
	"github.com/terdia/greenlight/infrastructures/persistence/postgres"
	"github.com/terdia/greenlight/internal/mailer"
	"github.com/terdia/greenlight/internal/registry"
)

var (
	buildTime string
	version   string
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

	cfg := new(config.Config)
	cfg.Version = version

	flag.IntVar(&cfg.AppPort, "port", 4000, "Api server")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.Db.Dsn, "dsn", "xxxxxx", "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.Smtp.Host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", "xxxxxxx", "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", "xxxxxxx", "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", "", "SMTP sender details")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})
	// Create a new version boolean flag with the default value of false.
	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)

		// Print out the contents of the buildTime variable.
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	logger := logger.New(os.Stdout, logger.LevelInfo)

	db, err := postgres.OpenDb(cfg.Db)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	mailer := mailer.New(cfg.Smtp)

	wg := new(sync.WaitGroup)

	expvar.NewString("appVersion").Set(cfg.Version)

	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

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
