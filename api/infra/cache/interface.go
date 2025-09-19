package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
	GetKeys(ctx context.Context, pattern string) ([]string, error)
	SetWithExpiry(ctx context.Context, key, value string, expiresAt time.Time) error
	SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error
	LoadScript(ctx context.Context, script string) (string, error)
	SetManyWithExpiry(ctx context.Context, items map[string]string, expiresAt time.Time) error
	ZIncrBy(ctx context.Context, key string, increment float64, member string) error
	ZAdd(ctx context.Context, key string, members map[string]float64) error
	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*ZRangeItem, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZRevRank(ctx context.Context, key string, member string) (int64, error)
	ZScore(ctx context.Context, key string, member string) (float64, error)
	Close()
}
