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
	Version     string      `env:"VERSION"     envDefault:"2.0.0"`

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
		Port             int    `env:"REST_PORT" envDefault:"8000"`
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
		StartTimeout  time.Duration `env:"INSTANCE_START_TIMEOUT" envDefault:"5m"`
		StopTimeout   time.Duration `env:"INSTANCE_STOP_TIMEOUT" envDefault:"2m"`
		ExtendTimeout time.Duration `env:"INSTANCE_EXTEND_TIMEOUT" envDefault:"5s"`
		Lifetime      time.Duration `env:"INSTANCE_LIFETIME" envDefault:"30m"`
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
