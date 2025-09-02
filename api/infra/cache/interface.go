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
	Close()
}
