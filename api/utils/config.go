package utils

import (
	"log"

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
