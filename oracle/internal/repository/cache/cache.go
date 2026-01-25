package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/oracle/infra/cache"
	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"

	"github.com/vmihailenco/msgpack/v5"
)

func CachedQuery[T any](ctx context.Context, cache cache.Cache, key string, ttl time.Duration, fetchFn func() (T, error)) (T, error) {
	var zero T

	// look in cache
	val, err := cache.Get(ctx, key)
	if err != nil {
		if !errorDom.Is(err, errorDom.ErrCacheMiss) {
			return zero, err
		}
	} else {
		return deserializeCache[T](ctx, val)
	}

	// fetch from the database
	result, err := fetchFn()
	if err != nil {
		return zero, err
	}

	// store in cache
	serialized, err := serializeCache(ctx, result)
	if err != nil {
		return zero, err
	}
	if err := cache.SetWithTTL(ctx, key, serialized, ttl); err != nil {
		return zero, err
	}

	return result, nil
}

func deserializeCache[T any](ctx context.Context, val string) (T, error) {
	var result T

	if err := msgpack.Unmarshal([]byte(val), &result); err != nil {
		return result, errorDom.Raise(
			ctx, errorDom.ErrUnmarshalError,
			"failed to unmarshal cache value", err,
			common.ExtraData{"type": fmt.Sprintf("%T", result), "value": val},
		)
	}
	return result, nil
}

func serializeCache[T any](ctx context.Context, val T) (string, error) {
	data, err := msgpack.Marshal(val)
	if err != nil {
		return "", errorDom.Raise(
			ctx, errorDom.ErrMarshalError,
			"failed to marshal cache value", err,
			common.ExtraData{"type": fmt.Sprintf("%T", val), "value": val},
		)
	}
	return string(data), nil
}
