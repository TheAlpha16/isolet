package logger

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	"github.com/TheAlpha16/isolet/api/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const loggerKey common.ContextKey = "logger"

var (
	logger *StandardLogger
	once   sync.Once
)

type StandardLogger struct {
	*zap.Logger
}

func Init() {
	once.Do(func() {
		logger = New()
	})
}

func New() *StandardLogger {
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

	newLogger := &StandardLogger{zapLogger}
	newLogger.Info("Logger initialized", zap.String("level", outputLevel.String()))

	return newLogger
}

func (l *StandardLogger) WithFields(fields ...zapcore.Field) *StandardLogger {
	return &StandardLogger{l.Logger.With(fields...)}
}

func (l *StandardLogger) Sync() {
	l.Logger.Sync()
}

func GetLogger(ctx context.Context) *StandardLogger {
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
	if logger == nil {
		Init()
	}
	return logger
}
