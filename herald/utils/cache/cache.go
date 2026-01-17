package cache

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/errors"

	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var globClient valkey.Client
var once sync.Once
var tracer = otel.Tracer("herald.cache")

func GetCache(ctx context.Context) (valkey.Client, error) {
	once.Do(func() {
		var err error
		globClient, err = createValkeyClient(ctx)
		if err != nil {
			panic(err)
		}
	})
	return globClient, nil
}

func createValkeyClient(ctx context.Context) (valkey.Client, error) {
	var client valkey.Client
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

	client, err = valkey.NewClient(opts)
	if err != nil {
		return nil, err
	}

	err = client.Do(ctx, client.B().Ping().Build()).Error()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Get(ctx context.Context, key string) (string, error) {
	ctx, span := tracer.Start(ctx, "get")
	defer span.End()

	client, err := GetCache(ctx)
	if err != nil {
		return "", err
	}

	span.SetAttributes(attribute.Bool("cache.hit", true))
	result, err := client.Do(ctx, client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			span.AddEvent("cache miss")
			span.SetAttributes(attribute.Bool("cache.hit", false))
			return "", nil
		}
		return "", errors.Raise(errors.ErrCacheCallFailed, "cache get failed", err)
	}

	return result, nil
}

func SetNXWithTTL(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	ctx, span := tracer.Start(ctx, "set_nx_with_ttl")
	defer span.End()

	client, err := GetCache(ctx)
	if err != nil {
		return false, err
	}

	result, err := client.Do(ctx, client.B().Set().Key(key).Value(value).Nx().Ex(ttl).Build()).AsBool()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return false, nil
		}
		return false, errors.Raise(errors.ErrCacheCallFailed, "cache set nx with ttl failed", err)
	}

	return result, nil
}
