package database

import (
	"context"
	"log"
	"time"

	"github.com/TheAlpha16/isolet/api/config"
)

func InitializeGlobalDBConfig() {
	RefreshCache()

	go func() {
		interval := config.DefaultConfigRefreshInterval

		for {
			time.Sleep(interval)

			RefreshCache()

			newInterval, err := config.GetDuration(config.ConfigRefreshIntervalKey)
			if err == nil {
				newInterval = newInterval * time.Second
				if newInterval != interval {
					interval = newInterval
				}
			} else {
				log.Println(err)
			}
		}
	}()
}

func RefreshCache() {
	ctx := context.Background()

	configs, err := GetConfigs(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	config.GlobalConfig.Mutex.Lock()
	defer config.GlobalConfig.Mutex.Unlock()

	for _, configEntry := range configs {
		config.GlobalConfig.Cache[configEntry.Key] = configEntry.Value
	}

	config.GlobalConfig.LastFetch = time.Now()
}
