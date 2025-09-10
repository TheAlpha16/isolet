package configvars

import (
	"context"
	"strconv"
	"sync"
	"time"

	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"go.uber.org/zap"
)

type cvImpl struct {
	repo  cvDom.Repository
	cache map[string]string
	mu    sync.RWMutex
}

func (cv *cvImpl) GetString(ctx context.Context, key cvDom.ConfigKey[string]) string {
	if val, ok := cv.get(key.Name); ok {
		return val
	}
	return key.Default
}

func (cv *cvImpl) GetBool(ctx context.Context, key cvDom.ConfigKey[bool]) bool {
	var err error
	var cast bool
	extraData := map[string]any{"key": key.Name, "expected": "bool"}

	if val, ok := cv.get(key.Name); ok {
		extraData["value"] = val
		if cast, err = strconv.ParseBool(val); err == nil {
			return cast
		}
	}

	cv.handleInvalid(ctx, key.Name, errorDom.Raise(ctx, errorDom.ErrConfigVarInvalid, "", err, extraData))
	return key.Default
}

func (cv *cvImpl) GetInt(ctx context.Context, key cvDom.ConfigKey[int]) int {
	var err error
	var cast int
	extraData := map[string]any{"key": key.Name, "expected": "int"}

	if val, ok := cv.get(key.Name); ok {
		extraData["value"] = val
		if cast, err = strconv.Atoi(val); err == nil {
			return cast
		}
	}

	cv.handleInvalid(ctx, key.Name, errorDom.Raise(ctx, errorDom.ErrConfigVarInvalid, "", err, extraData))
	return key.Default
}

func (cv *cvImpl) GetDuration(ctx context.Context, key cvDom.ConfigKey[time.Duration]) time.Duration {
	var err error
	var cast time.Duration
	extraData := map[string]any{"key": key.Name, "expected": "duration"}

	if val, ok := cv.get(key.Name); ok {
		extraData["value"] = val
		if cast, err = time.ParseDuration(val); err == nil {
			return cast
		}
	}

	cv.handleInvalid(ctx, key.Name, errorDom.Raise(ctx, errorDom.ErrConfigVarInvalid, "", err, extraData))
	return key.Default
}

func (cv *cvImpl) get(key string) (string, bool) {
	cv.mu.RLock()
	defer cv.mu.RUnlock()

	val, ok := cv.cache[key]
	return val, ok
}

func (cv *cvImpl) handleInvalid(ctx context.Context, key string, err error) {
	logger.GetAppLogger().Error("invalid config variable", zap.String("key", key), zap.Error(err))
	errorDom.RaiseToSentry(ctx, err)
}

func (cv *cvImpl) Refresh() {
	cache, err := cv.repo.Refresh(context.TODO())
	if err != nil {
		refreshErr := errorDom.Raise(context.TODO(), errorDom.ErrConfigVarRefreshFailed, "", err, nil)
		logger.GetAppLogger().Error("failed to refresh config variables", zap.Error(err))
		errorDom.RaiseToSentry(context.TODO(), refreshErr)
		return
	}

	cv.mu.Lock()
	defer cv.mu.Unlock()

	cv.cache = cache
}

func New(repo cvDom.Repository) cvDom.Usecase {
	cv := &cvImpl{
		repo:  repo,
		cache: make(map[string]string),
		mu:    sync.RWMutex{},
	}
	cv.Refresh()
	return cv
}
