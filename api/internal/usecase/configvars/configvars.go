package configvars

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"go.uber.org/zap"
)

var ErrMissingConfigVariable = errors.New("missing config variable")

type cvImpl struct {
	repo  cvDom.Repository
	cache map[string]string
	mu    sync.RWMutex
}

func (cv *cvImpl) GetString(key cvDom.ConfigKey[string]) string {
	if val, ok := cv.get(key.Name); ok {
		return val
	}
	return key.Default
}

func (cv *cvImpl) GetBool(key cvDom.ConfigKey[bool]) bool {
	var err error
	var cast bool

	if val, ok := cv.get(key.Name); ok {
		if cast, err = strconv.ParseBool(val); err == nil {
			return cast
		}
	}

	cv.handleInvalid(key.Name, err)
	return key.Default
}

func (cv *cvImpl) GetInt(key cvDom.ConfigKey[int]) int {
	var err error
	var cast int

	if val, ok := cv.get(key.Name); ok {
		if cast, err = strconv.Atoi(val); err == nil {
			return cast
		}
	}

	cv.handleInvalid(key.Name, err)
	return key.Default
}

func (cv *cvImpl) GetDuration(key cvDom.ConfigKey[time.Duration]) time.Duration {
	var err error
	var cast time.Duration

	if val, ok := cv.get(key.Name); ok {
		if cast, err = time.ParseDuration(val); err == nil {
			return cast
		}
	}

	cv.handleInvalid(key.Name, err)
	return key.Default
}

func (cv *cvImpl) get(key string) (string, bool) {
	cv.mu.RLock()
	defer cv.mu.RUnlock()

	if val, ok := cv.cache[key]; ok {
		return val, ok
	}

	cv.handleInvalid(key, ErrMissingConfigVariable)
	return "", false
}

func (cv *cvImpl) handleInvalid(key string, err error) {
	logger.GetAppLogger().Error("missing/invalid config variable", zap.String("key", key), zap.Error(err))
	errorDom.RaiseToSentry(context.TODO(), err)
}

func New(repo cvDom.Repository) *cvImpl {
	return &cvImpl{
		repo:  repo,
		cache: make(map[string]string),
		mu:    sync.RWMutex{},
	}
}
