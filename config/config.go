package config

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
