package config

import (
	"flag"
	"strings"
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
	Cors struct {
		TrustedOrigins []string
	}
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

	flag.Parse()

}

func NewConfig() *Config {
	return &cfg
}
