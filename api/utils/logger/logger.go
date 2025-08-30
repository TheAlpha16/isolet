package logger

import (
	"context"
	"fmt"

	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const loggerKey contextKey = "logger"

var logger *StandardLogger

type StandardLogger struct {
	*zap.Logger
}

func Init(config *utils.Config) {
	if logger != nil {
		return
	}

	var cfg zap.Config
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
}

func IntercepterLogger(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
	f := make([]zap.Field, 0, len(fields)/2)

	for i := 0; i < len(fields); i += 2 {
		key := fields[i]
		value := fields[i+1]

		switch v := value.(type) {
		case string:
			f = append(f, zap.String(key.(string), v))
		case int:
			f = append(f, zap.Int(key.(string), v))
		case bool:
			f = append(f, zap.Bool(key.(string), v))
		default:
			f = append(f, zap.Any(key.(string), v))
		}
	}

	logger := FromContext(ctx).With(f...)

	switch lvl {
	case logging.LevelDebug:
		logger.Debug(msg)
	case logging.LevelInfo:
		logger.Info(msg)
	case logging.LevelWarn:
		logger.Warn(msg)
	case logging.LevelError:
		logger.Error(msg)
	default:
		panic(fmt.Sprintf("unknown level %v", lvl))
	}
}

func (l *StandardLogger) WithFields(fields ...zapcore.Field) *StandardLogger {
	zapLogger := l.Logger
	for _, field := range fields {
		zapLogger = zapLogger.With(field)
	}
	return &StandardLogger{zapLogger}
}

func FromContext(ctx context.Context) *StandardLogger {
	if logger, ok := ctx.Value(loggerKey).(*StandardLogger); ok {
		return logger
	}
	return logger
}

func NewContext(
	ctx context.Context,
	l *StandardLogger,
) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}
