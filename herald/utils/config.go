package utils

import (
	"errors"
	"log"
	"os"
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
	Name        string      `env:"NAME" envDefault:"herald"`
	Environment Environment `env:"ENVIRONMENT" envDefault:"local"`
	Version     string      `env:"VERSION" envDefault:"0.1.0"`
	LogLevel    string      `env:"LOG_LEVEL" envDefault:"DEBUG"`
	Identity    string      `env:"IDENTITY" envDefault:"k8s_source"`

	Sentry struct {
		DSN             string  `env:"SENTRY_DSN" envDefault:""`
		SampleRate      float64 `env:"SENTRY_SAMPLE_RATE" envDefault:"1.0"`
		TraceSampleRate float64 `env:"SENTRY_TRACE_SAMPLE_RATE" envDefault:"0.05"`
	}

	K8s struct {
		KubeConfigFilePath string        `env:"K8S_KUBE_CONFIG_FILE_PATH" envDefault:""`
		Namespaces         []string      `env:"K8S_NAMESPACES" envDefault:"isolet,dynamic"`
		DeDupeWindow       time.Duration `env:"K8S_DEDUPE_WINDOW" envDefault:"60s"`
	}

	Valkey struct {
		Address  []string `env:"VALKEY_ADDRESS" envDefault:"valkey://localhost:6379"`
		Username string   `env:"VALKEY_USERNAME" envDefault:""`
		Password string   `env:"VALKEY_PASSWORD" envDefault:""`
		UseTLS   bool     `env:"VALKEY_USE_TLS" envDefault:"false"`
	}

	Emitter struct {
		Retries       int           `env:"EMITTER_RETRIES" envDefault:"5"`
		RetryInterval time.Duration `env:"EMITTER_RETRY_INTERVAL" envDefault:"1s"`
	}

	Kafka struct {
		Brokers string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
		Retries int    `env:"KAFKA_RETRIES" envDefault:"3"`
	}

	InstanceLifecycle struct {
		KafkaTopic      string `env:"INSTANCE_LIFECYCLE_KAFKA_TOPIC" envDefault:"herald.instance.lifecycle"`
		FactChannelSize int    `env:"INSTANCE_LIFECYCLE_FACT_CHANNEL_SIZE" envDefault:"1024"`
		Workers         int    `env:"INSTANCE_LIFECYCLE_WORKERS" envDefault:"10"`
	}

	Postgres struct {
		DSN            string        `env:"POSTGRES_DSN" envDefault:"postgres://postgres:postgres@localhost:5432/isolet?sslmode=disable"`
		ReplicationDSN string        `env:"POSTGRES_REPLICATION_DSN" envDefault:"postgres://postgres:postgres@localhost:5432/isolet?sslmode=disable&replication=database"`
		StandbyTimeout time.Duration `env:"POSTGRES_STANDBY_TIMEOUT" envDefault:"5s"`
	}

	Notification struct {
		KafkaTopic      string `env:"NOTIFICATION_KAFKA_TOPIC" envDefault:"herald.notifications"`
		FactChannelSize int    `env:"NOTIFICATION_FACT_CHANNEL_SIZE" envDefault:"1024"`
		Workers         int    `env:"NOTIFICATION_WORKERS" envDefault:"10"`
	}
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
