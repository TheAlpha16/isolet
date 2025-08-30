package logger

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/api/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const loggerKey contextKey = "logger"

var (
	logger *StandardLogger
	once   sync.Once
)

type StandardLogger struct {
	*zap.Logger
}

func Init() {
	once.Do(func() {
		var cfg zap.Config
		config := utils.GetConfig()
		outputLevel, err := zapcore.ParseLevel(config.LogLevel)
		if err != nil {
			outputLevel = zapcore.InfoLevel
		}

		if config.Environment != utils.LOCAL {
			cfg = zap.NewProductionConfig()
			cfg.Level = zap.NewAtomicLevelAt(outputLevel)
			cfg.Encoding = "json"
		} else {
			cfg = zap.NewDevelopmentConfig()
		}

		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stdout"}
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.TimeKey = "time"
		cfg.DisableStacktrace = true

		zapLogger, err := cfg.Build()
		if err != nil {
			panic(err)
		}

		logger = &StandardLogger{zapLogger}
	})
}

func (l *StandardLogger) WithFields(fields ...zapcore.Field) *StandardLogger {
	return &StandardLogger{l.Logger.With(fields...)}
}

func FromContext(ctx context.Context) *StandardLogger {
	if ctxLogger, ok := ctx.Value(loggerKey).(*StandardLogger); ok {
		return ctxLogger
	}
	return logger
}

func NewContext(
	ctx context.Context,
	l *StandardLogger,
) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func GetAppLogger() *StandardLogger {
	return logger
}
