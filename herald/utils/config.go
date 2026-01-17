package utils

import (
	"errors"
	"log"
	"os"

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
	Name        string      `env:"NAME" envDefault:"herald"`
	Environment Environment `env:"ENVIRONMENT" envDefault:"local"`
	Version     string      `env:"VERSION" envDefault:"0.1.0"`
	LogLevel    string      `env:"LOG_LEVEL" envDefault:"DEBUG"`
}

func GetConfig() *Config {
	if appConfig != nil {
		return appConfig
	} else {
		if err := godotenv.Load(); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Println(".env file not found, continuing without it")
			} else {
				log.Panic("Error loading .env file: ", err)
			}
		}

		cfg, err := env.ParseAs[Config]()
		if err != nil {
			log.Panic("Error parsing config: ", err)
		}
		return &cfg
	}
}
