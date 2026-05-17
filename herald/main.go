package main

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/herald/internal/app"
	"github.com/TheAlpha16/isolet/herald/internal/sources/k8s"
	"github.com/TheAlpha16/isolet/herald/internal/sources/postgres"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/cache"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/kafka"
	"github.com/TheAlpha16/isolet/herald/utils/logger"
	"github.com/TheAlpha16/isolet/herald/utils/tracer"

	"go.uber.org/zap"
)

func main() {
	globalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appLogger := logger.GetAppLogger()
	defer appLogger.Sync()

	wg := &sync.WaitGroup{}
	config := utils.GetConfig()
	tracer.InitSentry(config)
	if _, err := cache.GetCache(globalCtx); err != nil {
		appLogger.Fatal("failed to initialize cache", zap.Error(err))
	}

	kafkaClient, err := kafka.NewClient()
	if err != nil {
		errors.RaiseToSentry(globalCtx, err)
		appLogger.Fatal("failed to create kafka client", zap.Error(err))
	}

	utils.InterruptHandlerChannel <- func() {
		appLogger.Info("triggering context cancellation...")
		cancel()
		wg.Wait()
		appLogger.Info("all components stopped")
	}

	switch config.Identity {
	case utils.SourceK8s:
		app.RunPipeline(globalCtx, kafkaClient, wg, app.AppConfig{
			Source:      k8s.NewSource(globalCtx, wg),
			Workers:     config.InstanceLifecycle.Workers,
			ChannelSize: config.InstanceLifecycle.FactChannelSize,
			KafkaTopic:  config.InstanceLifecycle.KafkaTopic,
		})
	case utils.SourcePostgres:
		app.RunPipeline(globalCtx, kafkaClient, wg, app.AppConfig{
			Source:      postgres.NewSource(globalCtx, wg),
			Workers:     config.Notification.Workers,
			ChannelSize: config.Notification.FactChannelSize,
			KafkaTopic:  config.Notification.KafkaTopic,
		})
	}

	utils.InterruptHandler()
	appLogger.Info("herald shut down successfully")
}
