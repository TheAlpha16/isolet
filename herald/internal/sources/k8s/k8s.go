package k8s

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/internal/sources"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/tracer"

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

var k8sTracer = otel.Tracer("herald.sources.k8s")

type k8sSource struct {
	cache crcache.Cache
	wg    *sync.WaitGroup
}

func (s *k8sSource) Name() string {
	return utils.SourceK8s
}

func (s *k8sSource) Run(ctx context.Context, out chan<- facts.Fact) error {
	ctx, span, log := tracer.StartSpan(ctx, k8sTracer, "herald.sources.k8s.Run")
	defer span.End()

	if err := s.startInformers(ctx, out); err != nil {
		errors.HandleSpanError(ctx, span, log, "failed to start informers", err)
		return err
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		cacheCtx, cacheSpan, cacheLog := tracer.StartSpan(ctx, k8sTracer, "herald.sources.k8s.cache.Start")
		defer cacheSpan.End()

		if err := s.cache.Start(cacheCtx); err != nil {
			err = errors.Raise(errors.ErrK8sCacheFailed, "k8s cache failed to start", err)
			errors.HandleSpanError(cacheCtx, cacheSpan, cacheLog, "k8s cache failed to start", err)
			cacheLog.Fatal("failed to start cache", zap.Error(err))
		}
	}()

	if !s.cache.WaitForCacheSync(ctx) {
		err := errors.Raise(errors.ErrK8sCacheFailed, "timed out waiting for caches to sync", nil)
		errors.HandleSpanError(ctx, span, log, "cache sync timeout", err)
		return err
	}

	span.AddEvent("k8s.source.started")
	log.Info("K8s source started", zap.Strings("namespaces", utils.GetConfig().K8s.Namespaces))

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

func NewSource(ctx context.Context, wg *sync.WaitGroup) sources.Source {
	ctx, span, log := tracer.StartSpan(ctx, k8sTracer, "herald.sources.k8s.NewSource")
	defer span.End()

	cache, err := getK8sCache()
	if err != nil {
		errors.HandleSpanError(ctx, span, log, "failed to create k8s event cache", err)
		log.Fatal("failed to create k8s event cache", zap.Error(err))
	}

	return &k8sSource{
		cache: cache,
		wg:    wg,
	}
}
