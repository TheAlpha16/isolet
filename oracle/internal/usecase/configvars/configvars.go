package configvars

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/TheAlpha16/isolet/oracle/infra/cnc"
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"

	"go.uber.org/zap"
)

const (
	refreshCacheCmd = "refresh-cache"
)

type cvImpl struct {
	repo   cvDom.Repository
	cache  map[string]string
	cnc    cnc.CNC
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func (cv *cvImpl) GetString(ctx context.Context, key cvDom.ConfigKey[string]) string {
	if val, ok := cv.get(key.Name); ok {
		return val
	}
	return key.Default
}

func (cv *cvImpl) GetBool(ctx context.Context, key cvDom.ConfigKey[bool]) bool {
	return parseOrDefault(ctx, cv, key, strconv.ParseBool)
}

func (cv *cvImpl) GetInt(ctx context.Context, key cvDom.ConfigKey[int]) int {
	return parseOrDefault(ctx, cv, key, strconv.Atoi)
}

func (cv *cvImpl) GetDuration(ctx context.Context, key cvDom.ConfigKey[time.Duration]) time.Duration {
	return parseOrDefault(ctx, cv, key, time.ParseDuration)
}

func (cv *cvImpl) Refresh(ctx context.Context) error {
	if err := cv.cnc.Trigger(ctx, refreshCacheCmd, nil); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrConfigVarRefreshFailed, "", err, nil)
	}
	return nil
}

func (cv *cvImpl) refresh(ctx context.Context) {
	cache, err := cv.repo.Refresh(ctx)
	if err != nil {
		refreshErr := errorDom.Raise(ctx, errorDom.ErrConfigVarRefreshFailed, "", err, nil)
		logger.GetAppLogger().Error("failed to refresh config variables", zap.Error(err))
		errorDom.RaiseToSentry(ctx, refreshErr)
		return
	}

	cv.mu.Lock()
	defer cv.mu.Unlock()

	cv.cache = cache
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

func (cv *cvImpl) startAutoRefresh(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				cv.refresh(cv.ctx)
			case <-cv.ctx.Done():
				return
			}
		}
	}()
}

func parseOrDefault[T any](ctx context.Context, cv *cvImpl, key cvDom.ConfigKey[T], parser func(string) (T, error)) T {
	var val string
	var err error
	var cast T
	var ok bool
	extraData := map[string]any{
		"key":      key.Name,
		"expected": fmt.Sprintf("%T", key.Default),
	}

	if val, ok = cv.get(key.Name); !ok {
		return key.Default
	}

	extraData["value"] = val
	if cast, err = parser(val); err == nil {
		return cast
	}

	cv.handleInvalid(ctx, key.Name, errorDom.Raise(ctx, errorDom.ErrConfigVarInvalid, "", err, extraData))
	return key.Default
}

func New(ctx context.Context, repo cvDom.Repository, cncSvc cnc.CNC) cvDom.Usecase {
	config := utils.GetConfig()
	ctx, cancel := context.WithCancel(ctx)

	cv := &cvImpl{
		repo:   repo,
		cache:  make(map[string]string),
		cnc:    cncSvc,
		mu:     sync.RWMutex{},
		ctx:    ctx,
		cancel: cancel,
	}
	cv.refresh(ctx)
	cv.startAutoRefresh(config.ConfigVars.RefreshInterval)

	// register handler for distributed cache refresh
	if err := cv.cnc.Register(refreshCacheCmd, func(ctx context.Context, params map[string]any) error {
		cv.refresh(ctx)
		return nil
	}); err != nil {
		errorDom.RaiseToSentry(ctx, err)
		logger.GetAppLogger().Panic("failed to register config variable refresh handler", zap.Error(err))
	}

	utils.InterruptHandlerChannel <- func() {
		cancel()
	}

	return cv
}
