package tracer

import (
	"github.com/TheAlpha16/isolet/herald/utils"

	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func InitSentry(config *utils.Config) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.Sentry.DSN,
		AttachStacktrace: true,
		SampleRate:       config.Sentry.SampleRate,
		EnableTracing:    true,
		TracesSampleRate: config.Sentry.TraceSampleRate,
		ServerName:       config.Name,
		Release:          config.Version,
		Environment:      string(config.Environment),
	})
	if err != nil {
		panic(err)
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(config.Name)),
	)
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
		sdktrace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())
}
