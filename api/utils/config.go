package utils

import (
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environment string

var appConfig *Config

const (
	LOCAL Environment = "local"
	DEV   Environment = "dev"
	PROD  Environment = "prod"
)

type Config struct {
	Name        string      `env:"NAME"        envDefault:"api"`
	LogLevel    string      `env:"LOG_LEVEL"   envDefault:"DEBUG"`
	Environment Environment `env:"ENVIRONMENT" envDefault:"local"`
	Version     string      `env:"VERSION"     envDefault:"v2.0.0"`

	Database struct {
		Name                  string        `env:"DB_NAME"     envDefault:"isolet"`
		User                  string        `env:"DB_USER"     envDefault:"postgres"`
		Password              string        `env:"DB_PASSWORD" envDefault:"postgres"`
		Host                  string        `env:"DB_HOST"     envDefault:"localhost"`
		Port                  int           `env:"DB_PORT"     envDefault:"5432"`
		MaxConnections        int           `env:"DB_MAX_CONNECTIONS" envDefault:"50"`
		MaxConnectionIdleTime time.Duration `env:"DB_MAX_CONNECTION_IDLE_TIME" envDefault:"5m"`
	}

	Sentry struct {
		DSN             string  `env:"SENTRY_DSN" envDefault:""`
		SampleRate      float64 `env:"SENTRY_SAMPLE_RATE" envDefault:"1.0"`
		TraceSampleRate float64 `env:"SENTRY_TRACE_SAMPLE_RATE" envDefault:"1.0"`
	}
}

func GetConfig() *Config {
	if appConfig != nil {
		return appConfig
	} else {
		godotenv.Load()

		cfg, err := env.ParseAs[Config]()
		if err != nil {
			log.Panic("Error parsing config: ", err)
		}
		return &cfg
	}
}
