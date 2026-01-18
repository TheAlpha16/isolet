package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/TheAlpha16/isolet/herald/internal/emitter"
	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/logger"
	"go.uber.org/zap"
)

type pipeline struct {
	emitter emitter.Emitter
	workers int
	wg      *sync.WaitGroup
}

func (p *pipeline) Run(ctx context.Context, in <-chan facts.Fact) {
	log := logger.GetAppLogger()
	log.Info("starting pipeline", zap.Int("workers", p.workers))

	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func(workerID int) {
			defer p.wg.Done()
			p.runWorker(ctx, workerID, in)
		}(i)
	}

	utils.InterruptHandlerChannel <- func() {
		log.Info("pipeline shutting down")
		p.emitter.Close()
	}

	<-ctx.Done()
	log.Info("pipeline context canceled, waiting for workers...")
}

func (p *pipeline) runWorker(ctx context.Context, workerID int, in <-chan facts.Fact) {
	log := logger.GetAppLogger().With(zap.Int("worker", workerID))
	log.Debug("worker started")

	for {
		select {
		case <-ctx.Done():
			log.Debug("worker stopping due to context cancellation")
			return

		case fact, ok := <-in:
			if !ok {
				log.Debug("worker stopping because fact channel closed")
				return
			}
			p.emitWithRetry(ctx, log, fact)
		}
	}
}

func (p *pipeline) emitWithRetry(ctx context.Context, log *zap.Logger, fact facts.Fact) {
	for attempt := 1; attempt <= utils.GetConfig().Emitter.Retries; attempt++ {
		err := p.emitter.Emit(ctx, fact)
		if err == nil {
			log.Debug(
				"emitted fact successfully",
				zap.ByteString("key", fact.Key()),
				zap.String("fact_type", string(fact.FactType())),
			)
			return
		}

		log.Warn(
			"failed to emit fact",
			zap.ByteString("key", fact.Key()),
			zap.String("fact_type", string(fact.FactType())),
			zap.Int("attempt", attempt),
			zap.Error(err),
		)

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(attempt) * utils.GetConfig().Emitter.RetryInterval):
		}
	}

	// At-least-once semantics: last attempt
	log.Error(
		"giving up after retries, dropping fact",
		zap.ByteString("key", fact.Key()),
		zap.String("fact_type", string(fact.FactType())),
	)
}

func New(emitter emitter.Emitter, workers int, wg *sync.WaitGroup) *pipeline {
	if workers <= 0 {
		workers = 1
	}

	return &pipeline{
		emitter: emitter,
		workers: workers,
		wg:      wg,
	}
}
