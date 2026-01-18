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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type eventType string

const (
	eventTypeCreate eventType = "CREATE"
	eventTypeUpdate eventType = "UPDATE"
	eventTypeDelete eventType = "DELETE"
)

var objectHandlers = map[client.Object]func(obj any, out chan<- facts.Fact, eventType eventType){
	&tidev1.Instance{}: handleInstance,
}

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

	if err := s.startInformers(ctx, out); err != nil {
		errors.HandleSpanError(ctx, span, logger, "failed to start informers", err)
		return err
	}

	go func() {
		cacheCtx, cacheSpan := tracer.Start(ctx, "herald.sources.k8s.cache.Start")
		defer cacheSpan.End()

		if err := s.cache.Start(cacheCtx); err != nil {
			errors.HandleSpanError(cacheCtx, cacheSpan, logger, "k8s cache failed to start", err)
			logger.Fatal("failed to start cache", zap.Error(err))
		}
	}()

	if !s.cache.WaitForCacheSync(ctx) {
		err := errors.Raise(errors.ErrK8sCacheSyncFailed, "timed out waiting for caches to sync", nil)
		errors.HandleSpanError(ctx, span, logger, "cache sync timeout", err)
		return err
	}

	span.AddEvent("k8s.source.started")
	logger.Info("K8s source started", zap.Strings("namespaces", utils.GetConfig().K8s.Namespaces))

	<-ctx.Done()
	span.AddEvent("k8s.source.shutdown")
	return nil
}

func (s *k8sSource) startInformers(ctx context.Context, out chan<- facts.Fact) error {
	for objType, handlerFunc := range objectHandlers {
		informer, err := s.cache.GetInformer(ctx, objType)
		if err != nil {
			return errors.Raise(errors.ErrK8sInformerCreationFailed, "", err)
		}

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				handlerFunc(obj, out, eventTypeCreate)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				handlerFunc(newObj, out, eventTypeUpdate)
			},
			DeleteFunc: func(obj interface{}) {
				handlerFunc(obj, out, eventTypeDelete)
			},
		})
	}
	return nil
}

func NewSource(ctx context.Context) sources.Source {
	ctx, span := tracer.Start(ctx, "herald.sources.k8s.New")
	defer span.End()
	logger := logger.GetAppLogger()

	cache, err := getK8sCache()
	if err != nil {
		errors.HandleSpanError(ctx, span, logger, "failed to create k8s event cache", err)
		logger.Fatal("failed to create k8s event cache", zap.Error(err))
	}

	return &k8sSource{
		cache: cache,
	}
}
