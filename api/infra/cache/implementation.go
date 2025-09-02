package cache

import (
	"context"
	"strconv"
	"time"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/valkey-io/valkey-go"
	"go.opentelemetry.io/otel/attribute"
)

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	ctx, span := c.WithTrace(ctx, "get")
	defer span.End()

	span.SetAttributes(attribute.Bool("cache.hit", true))
	val, err := c.client.Do(ctx, c.client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		code := errorDom.ErrCacheCallFail
		if valkey.IsValkeyNil(err) {
			code = errorDom.ErrCacheMiss
			span.AddEvent("cache miss")
			span.SetAttributes(attribute.Bool("cache.hit", false))
		}
		return "", errorDom.Raise(ctx, code, "", err, nil)
	}
	return val, nil
}

func (c *cache) Set(ctx context.Context, key, value string) error {
	ctx, span := c.WithTrace(ctx, "set")
	defer span.End()

	err := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Build()).Error()
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	return nil
}

func (c *cache) Delete(ctx context.Context, key string) error {
	ctx, span := c.WithTrace(ctx, "delete")
	defer span.End()

	err := c.client.Do(ctx, c.client.B().Del().Key(key).Build()).Error()
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	return nil
}

func (c *cache) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	ctx, span := c.WithTrace(ctx, "keys")
	defer span.End()

	keys, err := c.client.Do(ctx, c.client.B().Keys().Pattern(pattern).Build()).AsStrSlice()
	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	span.SetAttributes(attribute.Int("cache.keys.count", len(keys)))
	return keys, nil
}

func (c *cache) SetWithExpiry(ctx context.Context, key, value string, expiresAt time.Time) error {
	ctx, span := c.WithTrace(ctx, "set_with_expiry")
	defer span.End()

	err := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Exat(expiresAt).Build()).Error()
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	return nil
}

func (c *cache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	ctx, span := c.WithTrace(ctx, "set_with_ttl")
	defer span.End()

	err := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Ex(ttl).Build()).Error()
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	return nil
}

func (c *cache) LoadScript(ctx context.Context, script string) (string, error) {
	ctx, span := c.WithTrace(ctx, "load_script")
	defer span.End()

	sha1Hash, err := c.client.Do(ctx, c.client.B().ScriptLoad().Script(script).Build()).ToString()
	if err != nil {
		return "", errorDom.Raise(ctx, errorDom.ErrCacheScriptLoadFail, "", err, nil)
	}
	return sha1Hash, nil
}

func (c *cache) SetManyWithExpiry(ctx context.Context, items map[string]string, expiresAt time.Time) error {
	ctx, span := c.WithTrace(ctx, "set_many")
	defer span.End()

	keys := make([]string, 0, len(items))
	args := make([]string, 0, len(items)+1)
	for k, v := range items {
		keys = append(keys, k)
		args = append(args, v)
	}
	args = append(args, strconv.FormatInt(expiresAt.Unix(), 10))

	err := c.client.Do(ctx, c.client.B().Evalsha().Sha1(setManyScriptHash).Numkeys(int64(len(keys))).Key(keys...).Arg(args...).Build()).Error()
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}

	return nil
}

func (c *cache) Close() {
	c.client.Close()
}
