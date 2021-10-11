package config

import (
	"flag"
	"os"
)

const version = "1.0.0"

var cfg Config

type Config struct {
	AppPort int
	Env     string
	Db      Db
	Version string
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	Smtp Smtp
}

type Db struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type Smtp struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

func init() {

	cfg.Version = version

	flag.IntVar(&cfg.AppPort, "port", 4000, "Api server")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.Db.Dsn, "dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.Smtp.Host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", os.Getenv("MAIL_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", os.Getenv("MAIL_PASSOWRD"), "SMTP password")
	flag.StringVar(&cfg.Smtp.Sender, "smtp-sender", "Terry Greenlight <no-reply@greenlight.com>", "SMTP sender")

	flag.Parse()

}

func NewConfig() *Config {
	return &cfg
}
