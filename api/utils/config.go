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
	LogLevel    string      `env:"LOG_LEVEL"   envDefault:"INFO"`
	Environment Environment `env:"ENVIRONMENT" envDefault:"local"`

	Database struct {
		User                  string        `env:"DB_USER"     envDefault:"user"`
		Password              string        `env:"DB_PASSWORD" envDefault:"password"`
		Host                  string        `env:"DB_HOST"     envDefault:"localhost"`
		Port                  int           `env:"DB_PORT"     envDefault:"5432"`
		Name                  string        `env:"DB_NAME"     envDefault:"isolet"`
		MaxConnections        int           `env:"DB_MAX_CONNECTIONS" envDefault:"50"`
		MaxConnectionIdleTime time.Duration `env:"DB_MAX_CONNECTION_IDLE_TIME" envDefault:"5m"`
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
