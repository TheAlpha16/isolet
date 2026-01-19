package configvars

import (
	"context"
	"time"
)

type Usecase interface {
	GetString(ctx context.Context, key ConfigKey[string]) string
	GetBool(ctx context.Context, key ConfigKey[bool]) bool
	GetInt(ctx context.Context, key ConfigKey[int]) int
	GetDuration(ctx context.Context, key ConfigKey[time.Duration]) time.Duration
	Refresh(ctx context.Context) error
}

type Repository interface {
	Refresh(ctx context.Context) (map[string]string, error)
}
