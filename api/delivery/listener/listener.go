package listener

import (
	"context"

	"github.com/TheAlpha16/isolet/api/infra"
	k8sInfra "github.com/TheAlpha16/isolet/api/infra/k8s"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	"github.com/TheAlpha16/isolet/api/internal/usecase"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"

	tideConstants "github.com/TheAlpha16/isolet/tide/utils"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	crcache "sigs.k8s.io/controller-runtime/pkg/cache"
)

func Start(ctx context.Context, usecases *usecase.Usecases, infra *infra.Infra) {
	log := logger.GetAppLogger()

	config, err := k8sInfra.GetRestConfig()
	if err != nil {
		log.Fatal("failed to get k8s config", zap.Error(err))
	}

	sch := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(sch))
	utilruntime.Must(tidev1.AddToScheme(sch))

	namespace := utils.GetConfig().Instances.Namespace

	k8sCache, err := crcache.New(config, crcache.Options{
		Scheme: sch,
		DefaultNamespaces: map[string]crcache.Config{
			namespace: {},
		},
	})
	if err != nil {
		log.Fatal("failed to create cache", zap.Error(err))
	}

	informer, err := k8sCache.GetInformer(ctx, &tidev1.Instance{})
	if err != nil {
		log.Fatal("failed to get informer", zap.Error(err))
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			instance, ok := obj.(*tidev1.Instance)
			if !ok {
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					log.Error("Couldn't get object from tombstone", zap.Any("obj", obj))
					return
				}
				instance, ok = tombstone.Obj.(*tidev1.Instance)
				if !ok {
					log.Error("Tombstone contained object that is not an Instance", zap.Any("obj", tombstone.Obj))
					return
				}
			}

			instDom := &instanceDom.Instance{
				ChallengeID: instance.Spec.Challenge.ID,
			}
			if instance.Spec.Team != nil {
				instDom.TeamID = &instance.Spec.Team.ID
			}

			if err := usecases.Instance.HandleEvent(context.Background(), tideConstants.EventReasonExpired, instDom); err != nil {
				log.Error("failed to handle instance expiry event", zap.Error(err))
			} else {
				log.Info("Handled instance deletion event", zap.String("instance", instance.Name))
			}
		},
	})

	go func() {
		if err := k8sCache.Start(ctx); err != nil {
			log.Fatal("failed to start cache", zap.Error(err))
		}
	}()

	if !k8sCache.WaitForCacheSync(ctx) {
		log.Log(zap.ErrorLevel, "Timed out waiting for caches to sync")
		return
	}

	log.Info("Listener started, watching for Instance events", zap.String("namespace", namespace))

	<-ctx.Done()
}
