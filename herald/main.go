package main

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/herald/internal/app"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/cache"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/kafka"
	"github.com/TheAlpha16/isolet/herald/utils/logger"

	"go.uber.org/zap"
)

func main() {
	globalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appLogger := logger.GetAppLogger()
	defer appLogger.Sync()

	wg := &sync.WaitGroup{}
	config := utils.GetConfig()
	utils.InitSentry(config)
	cache.GetCache(globalCtx)

	kafkaClient, err := kafka.NewClient()
	if err != nil {
		errors.RaiseToSentry(globalCtx, err)
		appLogger.Fatal("failed to create kafka client", zap.Error(err))
	}

	utils.InterruptHandlerChannel <- func() {
		appLogger.Info("Triggering context cancellation...")
		cancel()
		wg.Wait()
		appLogger.Info("All components stopped.")
	}

	app.Start(globalCtx, kafkaClient, wg)
}
