package app

import (
	"context"
	"sync"

	kafkaEmit "github.com/TheAlpha16/isolet/herald/internal/emitter/kafka"
	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/internal/pipeline"
	"github.com/TheAlpha16/isolet/herald/internal/sources"
	"github.com/TheAlpha16/isolet/herald/internal/sources/k8s"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/kafka"
)

type AppConfig struct {
	Source      sources.Source
	Workers     int
	ChannelSize int
	KafkaTopic  string
}

func Start(ctx context.Context, kafkaClient *kafka.Client, wg *sync.WaitGroup) {
	config := utils.GetConfig()

	RunPipeline(ctx, kafkaClient, wg, AppConfig{
		Source:      k8s.NewSource(ctx),
		Workers:     config.InstanceLifecycle.Workers,
		ChannelSize: config.InstanceLifecycle.FactChannelSize,
		KafkaTopic:  config.InstanceLifecycle.KafkaTopic,
	})

	utils.InterruptHandler()
}

func RunPipeline(ctx context.Context, kafkaClient *kafka.Client, wg *sync.WaitGroup, config AppConfig) {
	factChannel := make(chan facts.Fact, config.ChannelSize)

	k8sSource := k8s.NewSource(ctx)
	emitter := kafkaEmit.NewEmitter(config.KafkaTopic, kafkaClient)
	pipeline := pipeline.New(emitter, config.Workers, wg)
	go pipeline.Run(ctx, factChannel)
	go k8sSource.Run(ctx, factChannel)
}
