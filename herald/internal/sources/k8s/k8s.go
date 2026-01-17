package k8s

import (
	"context"

	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/internal/sources"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/logger"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
	crcache "sigs.k8s.io/controller-runtime/pkg/cache"
)

var tracer = otel.Tracer("herald.sources.k8s")

type k8sSource struct {
	cache crcache.Cache
}

func (s *k8sSource) Name() string {
	return utils.SourceK8s
}

func (s *k8sSource) Run(ctx context.Context, out chan<- facts.Fact) error {
	ctx, span := tracer.Start(ctx, "herald.sources.k8s.Run")
	defer span.End()
	logger := logger.GetAppLogger()

	informer, err := s.cache.GetInformer(ctx, &tidev1.Instance{})
	if err != nil {
		handleError(ctx, span, logger, "failed to get informer", err, nil)
		return err
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{})

	if !s.cache.WaitForCacheSync(ctx) {
		err := errors.Raise(errors.ErrInternalError, "timed out waiting for caches to sync", nil)
		handleError(ctx, span, logger, "cache sync timeout", err, nil)
		return err
	}

	span.AddEvent("k8s.source.started")
	logger.Info("K8s source started", zap.Strings("namespaces", utils.GetConfig().K8s.Namespaces))

	<-ctx.Done()
	span.AddEvent("k8s.source.shutdown")
	return nil
}

func New(ctx context.Context) sources.Source {
	ctx, span := tracer.Start(ctx, "herald.sources.k8s.New")
	defer span.End()
	logger := logger.GetAppLogger()

	cache, err := getK8sCache()
	if err != nil {
		handleError(ctx, span, logger, "failed to create k8s event cache", err, nil)
		logger.Fatal("failed to create k8s event cache", zap.Error(err))
	}

	return &k8sSource{
		cache: cache,
	}
}
