package app

import (
	"context"
	"sync"

	kafkaEmit "github.com/TheAlpha16/isolet/herald/internal/emitter/kafka"
	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/internal/pipeline"
	"github.com/TheAlpha16/isolet/herald/internal/sources"
	"github.com/TheAlpha16/isolet/herald/utils/kafka"
	"github.com/TheAlpha16/isolet/herald/utils/logger"

	"go.uber.org/zap"
)

type AppConfig struct {
	Source      sources.Source
	Workers     int
	ChannelSize int
	KafkaTopic  string
}

func RunPipeline(ctx context.Context, kafkaClient *kafka.Client, wg *sync.WaitGroup, config AppConfig) {
	factChannel := make(chan facts.Fact, config.ChannelSize)
	emitter := kafkaEmit.NewEmitter(config.KafkaTopic, kafkaClient)
	pipeline := pipeline.New(ctx, emitter, config.Workers, wg)

	// start the pipeline
	wg.Add(1)
	go func() {
		defer wg.Done()
		pipeline.Run(ctx, factChannel)
	}()

	// start the source
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := config.Source.Run(ctx, factChannel); err != nil {
			logger.GetAppLogger().Fatal("error from source", zap.Error(err))
		}
	}()
}
