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
		Name                  string        `env:"DB_NAME" envDefault:"isolet"`
		User                  string        `env:"DB_USER" envDefault:"postgres"`
		Password              string        `env:"DB_PASSWORD" envDefault:"postgres"`
		Host                  string        `env:"DB_HOST" envDefault:"localhost"`
		Port                  int           `env:"DB_PORT" envDefault:"5432"`
		MaxOpenConnections    int           `env:"DB_MAX_OPEN_CONNECTIONS" envDefault:"30"`
		MaxIdleConnections    int           `env:"DB_MAX_IDLE_CONNECTIONS" envDefault:"30"`
		MaxConnectionLifeTime time.Duration `env:"DB_MAX_CONNECTION_LIFETIME" envDefault:"3600s"`
	}

	Sentry struct {
		DSN             string  `env:"SENTRY_DSN" envDefault:""`
		SampleRate      float64 `env:"SENTRY_SAMPLE_RATE" envDefault:"1.0"`
		TraceSampleRate float64 `env:"SENTRY_TRACE_SAMPLE_RATE" envDefault:"0.05"`
	}

	Rest struct {
		Port int `env:"REST_PORT" envDefault:"8000"`
	}

	Valkey struct {
		Address  []string `env:"VALKEY_ADDRESS" envDefault:"valkey://localhost:6379"`
		Username string   `env:"VALKEY_USERNAME" envDefault:""`
		Password string   `env:"VALKEY_PASSWORD" envDefault:""`
		UseTLS   bool     `env:"VALKEY_USE_TLS" envDefault:"false"`
	}

	Token struct {
		EmailVerificationValidity time.Duration `env:"EMAIL_VERIFICATION_VALIDITY" envDefault:"600s"`
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
