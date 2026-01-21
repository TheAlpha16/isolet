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
type Identity string

var appConfig *Config

const (
	LOCAL Environment = "local"
	DEV   Environment = "dev"
	PROD  Environment = "prod"

	IdentityRest     Identity = "rest"
	IdentityConsumer Identity = "consumer"
)

type Config struct {
	Name        string      `env:"NAME"        envDefault:"oracle"`
	LogLevel    string      `env:"LOG_LEVEL"   envDefault:"DEBUG"`
	Environment Environment `env:"ENVIRONMENT" envDefault:"local"`
	Version     string      `env:"VERSION"     envDefault:"2.0.0"`
	Identity    Identity    `env:"IDENTITY"    envDefault:"rest"`

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
		Port             int    `env:"REST_PORT" envDefault:"80"`
		APIVersionPrefix string `env:"REST_API_VERSION_PREFIX" envDefault:"/api/v1"`
	}

	Valkey struct {
		Address  []string `env:"VALKEY_ADDRESS" envDefault:"valkey://localhost:6379"`
		Username string   `env:"VALKEY_USERNAME" envDefault:""`
		Password string   `env:"VALKEY_PASSWORD" envDefault:""`
		UseTLS   bool     `env:"VALKEY_USE_TLS" envDefault:"false"`
	}

	Token struct {
		SigningKey                string        `env:"TOKEN_SIGNING_KEY" envDefault:"trustmebro"`
		EmailVerificationValidity time.Duration `env:"TOKEN_EMAIL_VERIFICATION_VALIDITY" envDefault:"10m"`
		AuthValidity              time.Duration `env:"TOKEN_AUTH_VALIDITY" envDefault:"24h"`
		PasswordResetValidity     time.Duration `env:"TOKEN_PASSWORD_RESET_VALIDITY" envDefault:"30m"`
		TeamInviteValidity        time.Duration `env:"TOKEN_TEAM_INVITE_VALIDITY" envDefault:"1h"`
	}

	CNC struct {
		Channel string `env:"CNC_CHANNEL" envDefault:"cnc"`
	}

	ConfigVars struct {
		RefreshInterval time.Duration `env:"CONFIG_VARS_REFRESH_INTERVAL" envDefault:"10m"`
	}

	SMTP struct {
		Retries     int           `env:"SMTP_RETRIES" envDefault:"3"`
		Timeout     time.Duration `env:"SMTP_TIMEOUT" envDefault:"10s"`
		ChannelSize int           `env:"SMTP_CHANNEL_SIZE" envDefault:"10"`
	}

	Challenges struct {
		CacheTTL time.Duration `env:"CHALLENGE_CACHE_TTL" envDefault:"1h"`
	}

	Hints struct {
		CacheTTL time.Duration `env:"HINT_CACHE_TTL" envDefault:"1h"`
	}

	Instances struct {
		Namespace           string        `env:"INSTANCE_NAMESPACE" envDefault:"isolet"`
		StartTimeout        time.Duration `env:"INSTANCE_START_TIMEOUT" envDefault:"5m"`
		StopTimeout         time.Duration `env:"INSTANCE_STOP_TIMEOUT" envDefault:"2m"`
		ExtendTimeout       time.Duration `env:"INSTANCE_EXTEND_TIMEOUT" envDefault:"5s"`
		Lifetime            time.Duration `env:"INSTANCE_LIFETIME" envDefault:"30m"`
		MaxLifetime         time.Duration `env:"INSTANCE_MAX_LIFETIME" envDefault:"2h"`
		LimitCPU            string        `env:"INSTANCE_LIMIT_CPU" envDefault:"50m"`
		LimitMemory         string        `env:"INSTANCE_LIMIT_MEMORY" envDefault:"128Mi"`
		SecretKey           string        `env:"INSTANCE_SECRET_KEY" envDefault:"trustmebro"`
		FactTopic           string        `env:"INSTANCE_FACT_TOPIC" envDefault:"herald.instance.lifecycle"`
		ExpiryCheckInterval time.Duration `env:"INSTANCE_EXPIRY_CHECK_INTERVAL" envDefault:"5m"`
	}

	K8s struct {
		KubeConfigFilePath string        `env:"K8S_KUBE_CONFIG_FILE_PATH" envDefault:""`
		InstanceKind       string        `env:"K8S_INSTANCE_KIND" envDefault:"Instance"`
		InstanceAPIVersion string        `env:"K8S_INSTANCE_API_VERSION" envDefault:"challenges.isolet.dev/v1"`
		InstancePollRate   time.Duration `env:"K8S_INSTANCE_POLL_RATE" envDefault:"5s"`
		InsecureSkipVerify bool          `env:"K8S_INSECURE_SKIP_VERIFY" envDefault:"false"`
	}

	Manifest struct {
		CacheTTL time.Duration `env:"MANIFEST_CACHE_TTL" envDefault:"1h"`
	}

	Kafka struct {
		Brokers string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
		GroupID string `env:"KAFKA_GROUP_ID" envDefault:"oracle-consumer"`
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
