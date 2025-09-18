package cache

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const setManyScript = `
	--[[
		This script sets multiple key-value pairs atomically.
		Each key uses the same expiry timestamp passed as the last argument in ARGV.
		Usage pattern:
		  KEYS = list of keys
		  ARGV = [value1, value2, ..., expiryTimestamp]
	--]]
	for i = 1, #KEYS do
		redis.call("SET", KEYS[i], ARGV[i], "EXAT", ARGV[#ARGV])
	end
	return true
`

var (
	setManyScriptHash string
)

type cache struct {
	client valkey.Client
	tracer trace.Tracer
}

type ZRangeItem struct {
	Member string
	Score  float64
}

func (c *cache) WithTrace(ctx context.Context, operation string) (context.Context, trace.Span) {
	return c.tracer.Start(
		ctx,
		fmt.Sprintf("valkey.%s", operation),
		trace.WithAttributes(attribute.String("cache.operation", operation)),
	)
}

func NewClient(ctx context.Context, client valkey.Client) (Cache, error) {
	var err error

	tracer := otel.GetTracerProvider().Tracer("api.cache", trace.WithInstrumentationAttributes(attribute.String("cache.provider", "valkey")))

	cache := &cache{
		client: client,
		tracer: tracer,
	}

	setManyScriptHash, err = cache.LoadScript(ctx, setManyScript)
	if err != nil {
		return nil, err
	}

	return cache, nil
}
