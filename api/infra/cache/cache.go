package cache

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type cache struct {
	client valkey.Client
	tracer trace.Tracer
}

func (c *cache) WithTrace(ctx context.Context, operation string) (context.Context, trace.Span) {
	return c.tracer.Start(
		ctx,
		fmt.Sprintf("valkey.%s", operation),
		trace.WithAttributes(attribute.String("cache.operation", operation)),
	)
}

func NewClient(ctx context.Context) (Cache, error) {
	config := utils.GetConfig()
	opts, err := valkey.ParseURL(config.Valkey.Address[0])
	if err != nil {
		return nil, err
	}
	opts.ClientName = config.Name
	opts.Username = config.Valkey.Username
	opts.Password = config.Valkey.Password

	if config.Valkey.UseTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, err
	}

	err = client.Do(ctx, client.B().Ping().Build()).Error()
	if err != nil {
		return nil, err
	}

	tracer := otel.GetTracerProvider().Tracer("api.cache", trace.WithInstrumentationAttributes(attribute.String("cache.provider", "valkey")))

	return &cache{
		client: client,
		tracer: tracer,
	}, nil
}
