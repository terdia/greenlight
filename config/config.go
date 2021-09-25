package config

import (
	"flag"
	"os"
)

const version = "1.0.0"

var cfg Config
var db Db

type Config struct {
	AppPort int
	Env     string
	Db      Db
	Version string
}

type Db struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func init() {
	cfg.Version = version

	flag.IntVar(&cfg.AppPort, "port", 4000, "Api server")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&db.Dsn, "dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	cfg.Db = db

	flag.Parse()

}

func NewConfig() *Config {
	return &cfg
}
