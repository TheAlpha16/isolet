package listener

import (
	"context"

	"github.com/TheAlpha16/isolet/api/infra"
	k8sInfra "github.com/TheAlpha16/isolet/api/infra/k8s"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/internal/usecase"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	crcache "sigs.k8s.io/controller-runtime/pkg/cache"
)

var tracer = otel.Tracer("api/delivery/listener")

func Start(ctx context.Context, usecases *usecase.Usecases, infra *infra.Infra) {
	ctx, span := tracer.Start(ctx, "listener.Start")
	defer span.End()

	logger := logger.GetAppLogger()

	config, err := k8sInfra.GetRestConfig()
	if err != nil {
		handleError(ctx, span, logger, "failed to get k8s config", err, nil)
		logger.Fatal("failed to get k8s config", zap.Error(err))
	}

	sch := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(sch))
	utilruntime.Must(tidev1.AddToScheme(sch))

	namespace := utils.GetConfig().Instances.Namespace
	span.SetAttributes(attribute.String("k8s.namespace", namespace))

	k8sCache, err := crcache.New(config, crcache.Options{
		Scheme: sch,
		DefaultNamespaces: map[string]crcache.Config{
			namespace: {},
		},
	})
	if err != nil {
		handleError(ctx, span, logger, "failed to create cache", err, nil)
		logger.Fatal("failed to create cache", zap.Error(err))
	}

	informer, err := k8sCache.GetInformer(ctx, &tidev1.Instance{})
	if err != nil {
		handleError(ctx, span, logger, "failed to get informer", err, nil)
		logger.Fatal("failed to get informer", zap.Error(err))
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			handleDeleteEvent(obj, usecases)
		},
	})

	go func() {
		cacheCtx, cacheSpan := tracer.Start(ctx, "listener.k8sCache.Start")
		defer cacheSpan.End()

		if err := k8sCache.Start(cacheCtx); err != nil {
			handleError(cacheCtx, cacheSpan, logger, "cache failed to start", err, nil)
			logger.Fatal("failed to start cache", zap.Error(err))
		}
	}()

	if !k8sCache.WaitForCacheSync(ctx) {
		err := errorDom.RaiseInternal(ctx, "timed out waiting for caches to sync", nil, nil)
		handleError(ctx, span, logger, "cache sync timeout", err, nil)
		return
	}

	span.AddEvent("listener.started")
	logger.Info("Listener started, watching for Instance events", zap.String("namespace", namespace))

	<-ctx.Done()
	span.AddEvent("listener.shutdown")
}
