package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
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

	if err := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Build()).Error(); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	return nil
}

func (c *cache) Delete(ctx context.Context, key string) error {
	ctx, span := c.WithTrace(ctx, "delete")
	defer span.End()

	if err := c.client.Do(ctx, c.client.B().Del().Key(key).Build()).Error(); err != nil {
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

	if err := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Exat(expiresAt).Build()).Error(); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	return nil
}

func (c *cache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	ctx, span := c.WithTrace(ctx, "set_with_ttl")
	defer span.End()

	if err := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Ex(ttl).Build()).Error(); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	return nil
}

func (c *cache) SetNXWithTTL(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	ctx, span := c.WithTrace(ctx, "set_nx_with_ttl")
	defer span.End()

	set, err := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Nx().Ex(ttl).Build()).AsBool()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return false, nil
		}
		return false, errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}
	return set, nil
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

	if err := c.client.Do(ctx, c.client.B().Evalsha().Sha1(setManyScriptHash).Numkeys(int64(len(keys))).Key(keys...).Arg(args...).Build()).Error(); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}

	return nil
}

func (c *cache) ZIncrBy(ctx context.Context, key string, increment float64, member string) error {
	ctx, span := c.WithTrace(ctx, "zincr")
	defer span.End()

	if err := c.client.Do(ctx, c.client.B().Zincrby().Key(key).Increment(increment).Member(member).Build()).Error(); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}

	return nil
}

func (c *cache) ZAdd(ctx context.Context, key string, members map[string]float64) error {
	ctx, span := c.WithTrace(ctx, "zadd")
	defer span.End()

	query := c.client.B().Zadd().Key(key).ScoreMember()
	for member, score := range members {
		query = query.ScoreMember(score, member)
	}

	if err := c.client.Do(ctx, query.Build()).Error(); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}

	return nil
}

func (c *cache) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*ZRangeItem, error) {
	ctx, span := c.WithTrace(ctx, "zrevrange_with_scores")
	defer span.End()

	result, err := c.client.Do(ctx, c.client.B().Zrevrange().Key(key).Start(start).Stop(stop).Withscores().Build()).AsZScores()
	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}

	items := make([]*ZRangeItem, 0, len(result))
	for _, item := range result {
		items = append(items, &ZRangeItem{
			Member: item.Member,
			Score:  item.Score,
		})
	}

	return items, nil
}

func (c *cache) ZCard(ctx context.Context, key string) (int64, error) {
	ctx, span := c.WithTrace(ctx, "zcard")
	defer span.End()

	result, err := c.client.Do(ctx, c.client.B().Zcard().Key(key).Build()).AsInt64()
	if err != nil {
		return 0, errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}

	return result, nil
}

func (c *cache) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	ctx, span := c.WithTrace(ctx, "zrevrank")
	defer span.End()

	result, err := c.client.Do(ctx, c.client.B().Zrevrank().Key(key).Member(member).Build()).AsInt64()
	if err != nil {
		if err == valkey.Nil {
			return 0, errorDom.Raise(ctx, errorDom.ErrCacheZSetMissingMember, "", err, common.ExtraData{"key": key, "member": member})
		}
		return 0, errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}

	return result, nil
}

func (c *cache) ZScore(ctx context.Context, key string, member string) (float64, error) {
	ctx, span := c.WithTrace(ctx, "zscore")
	defer span.End()

	result, err := c.client.Do(ctx, c.client.B().Zscore().Key(key).Member(member).Build()).AsFloat64()
	if err != nil {
		if err == valkey.Nil {
			return 0, errorDom.Raise(ctx, errorDom.ErrCacheZSetMissingMember, "", err, common.ExtraData{"key": key, "member": member})
		}
		return 0, errorDom.Raise(ctx, errorDom.ErrCacheCallFail, "", err, nil)
	}

	return result, nil
}

func (c *cache) Close() {
	c.client.Close()
}
