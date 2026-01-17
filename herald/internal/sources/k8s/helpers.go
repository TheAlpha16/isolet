package k8s

import (
	"context"
	"strings"

	"github.com/TheAlpha16/isolet/herald/utils"
	errorDom "github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/k8s"
	"github.com/TheAlpha16/isolet/herald/utils/logger"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	crcache "sigs.k8s.io/controller-runtime/pkg/cache"
)

func handleError(ctx context.Context, span trace.Span, log *logger.StandardLogger, msg string, err error, extraFields map[string]any) {
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)
	errorDom.RaiseToSentry(ctx, err)

	zapFields := []zap.Field{zap.Error(err)}
	for k, v := range extraFields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	log.Error(msg, zapFields...)
}

func getK8sCache() (crcache.Cache, error) {
	config, err := k8s.GetRestConfig()
	if err != nil {
		return nil, err
	}

	sch := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(sch))
	utilruntime.Must(tidev1.AddToScheme(sch))

	namespaces := utils.GetConfig().K8s.Namespaces

	namespaceMap := make(map[string]crcache.Config)
	for _, ns := range namespaces {
		namespaceMap[ns] = crcache.Config{}
	}

	cache, err := crcache.New(config, crcache.Options{
		Scheme:            sch,
		DefaultNamespaces: namespaceMap,
	})
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func extractObject[T any](obj any) (T, bool) {
	var zero T

	if castedObj, ok := obj.(T); ok {
		return castedObj, true
	}

	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if !ok {
		return zero, false
	}

	castedObj, ok := tombstone.Obj.(T)
	return castedObj, ok
}

func getCacheKey(parts ...string) string {
	return "herald:k8s:" + strings.Join(parts, ":")
}
