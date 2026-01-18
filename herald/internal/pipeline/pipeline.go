package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/TheAlpha16/isolet/herald/internal/emitter"
	"github.com/TheAlpha16/isolet/herald/pkg/facts"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/logger"
	"github.com/TheAlpha16/isolet/herald/utils/tracer"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var pipelineTracer = otel.Tracer("herald.pipeline")

type pipeline struct {
	emitter emitter.Emitter
	workers int
	wg      *sync.WaitGroup
}

func (p *pipeline) Run(ctx context.Context, in <-chan facts.Fact) {
	ctx, span, log := tracer.StartSpan(ctx, pipelineTracer, "herald.pipeline.Run")
	defer span.End()

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
	span.AddEvent("pipeline.shutdown")
}

func (p *pipeline) runWorker(ctx context.Context, workerID int, in <-chan facts.Fact) {
	ctx, span, log := tracer.StartSpan(ctx, pipelineTracer, "herald.pipeline.runWorker")
	defer span.End()

	log = log.With(zap.Int("worker", workerID))
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
	var err error
	ctx, span, log := tracer.StartSpan(ctx, pipelineTracer, "herald.pipeline.emitWithRetry")
	defer span.End()

	log = log.With(
		zap.ByteString("key", fact.Key()),
		zap.String("fact_type", string(fact.FactType())),
	)
	log.Debug("received fact in pipeline")

	for attempt := 1; attempt <= utils.GetConfig().Emitter.Retries; attempt++ {
		err = p.emitter.Emit(ctx, fact)
		if err == nil {
			log.Debug("emitted fact successfully")
			return
		}

		log.Error(
			"failed to emit fact",
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
	errors.HandleSpanError(ctx, span, log, "giving up after retries, dropping fact", err)
}

func New(emitter emitter.Emitter, workers int, wg *sync.WaitGroup) *pipeline {
	if workers <= 0 {
		logger.GetAppLogger().Warn("invalid worker count, defaulting to 1")
		workers = 1
	}

	return &pipeline{
		emitter: emitter,
		workers: workers,
		wg:      wg,
	}
}
