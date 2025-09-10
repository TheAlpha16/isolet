package configvars

import (
	"context"
	"fmt"
	"sync"
	"time"

	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"go.uber.org/zap"
)

type cvImpl struct {
	repo  cvDom.Repository
	cache map[string]any
	mu    sync.RWMutex
}

func (cv *cvImpl) GetString(key cvDom.ConfigKey[string]) string {
	cv.mu.RLock()
	defer cv.mu.RUnlock()

	if v, ok := cv.cache[key.Name]; ok {
		if cast, ok := v.(string); ok {
			return cast
		}
	}
	cv.handleInvalid(key.Name)
	return key.Default
}

func (cv *cvImpl) GetBool(key cvDom.ConfigKey[bool]) bool {
	cv.mu.RLock()
	defer cv.mu.RUnlock()

	if v, ok := cv.cache[key.Name]; ok {
		if cast, ok := v.(bool); ok {
			return cast
		}
	}
	cv.handleInvalid(key.Name)
	return key.Default
}

func (cv *cvImpl) GetInt(key cvDom.ConfigKey[int]) int {
	cv.mu.RLock()
	defer cv.mu.RUnlock()

	if v, ok := cv.cache[key.Name]; ok {
		if cast, ok := v.(int); ok {
			return cast
		}
	}
	cv.handleInvalid(key.Name)
	return key.Default
}

func (cv *cvImpl) GetDuration(key cvDom.ConfigKey[time.Duration]) time.Duration {
	cv.mu.RLock()
	defer cv.mu.RUnlock()

	if v, ok := cv.cache[key.Name]; ok {
		if cast, ok := v.(time.Duration); ok {
			return cast
		}
	}
	cv.handleInvalid(key.Name)
	return key.Default
}

func (cv *cvImpl) handleInvalid(key string) {
	logger.GetAppLogger().Error("missing/invalid config variable", zap.String("key", key))
	errorDom.RaiseToSentry(context.TODO(), fmt.Errorf("invalid config value for key %s", key))
}

func New(repo cvDom.Repository) *cvImpl {
	return &cvImpl{
		repo:  repo,
		cache: make(map[string]any),
		mu:    sync.RWMutex{},
	}
}
