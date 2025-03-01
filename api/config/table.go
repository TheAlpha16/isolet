package config

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	ErrConfigNotFound = errors.New("configuration key not found")
	ErrParsingValue   = errors.New("error parsing configuration value")
)

type ConfigCache struct {
	Cache     map[string]string
	Mutex     sync.RWMutex
	LastFetch time.Time
}

var (
	GlobalConfig = &ConfigCache{
		Cache: make(map[string]string),
	}
	DefaultConfigRefreshInterval = 5 * time.Second
	ConfigRefreshIntervalKey = "SYNC_CONFIG_SECONDS"
)

func Get(key string) (string, error) {
	GlobalConfig.Mutex.RLock()
	defer GlobalConfig.Mutex.RUnlock()

	value, exists := GlobalConfig.Cache[key]
	if !exists {
		return "", fmt.Errorf("%w: %s", ErrConfigNotFound, key)
	}
	return value, nil
}

func GetInt(key string) (int64, error) {
	val, err := Get(key)
	if err != nil {
		return 0, err
	}
	
	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s (%s is not a valid integer)", ErrParsingValue, key, val)
	}
	return intVal, nil
}

func GetBool(key string) (bool, error) {
	val, err := Get(key)
	if err != nil {
		return false, err
	}
	
	if val == "true" {
		return true, nil
	} else if val == "false" {
		return false, nil
	} else {
		return false, fmt.Errorf("%w: %s (%s is not a valid boolean)", ErrParsingValue, key, val)
	}
}

func GetDuration(key string) (time.Duration, error) {
	val, err := GetInt(key)
	if err != nil {
		return 0, err
	}
	return time.Duration(val)* time.Second, nil
}
