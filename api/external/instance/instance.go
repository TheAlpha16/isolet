package instance

import (
	"context"
	"net/http"

	k8sInfra "github.com/TheAlpha16/isolet/api/infra/k8s"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
	"github.com/TheAlpha16/isolet/api/utils"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"github.com/TheAlpha16/isolet/tide/sdk"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type instanceSvc struct {
	client sdk.Handler
}

func (is *instanceSvc) Start(ctx context.Context, inst *instanceDom.Instance) error {
	instance := toTideInstance(ctx, inst)

	if err := is.client.CreateInstance(ctx, instance); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceCreationFailed, "", err, nil)
	}

	if err := is.reconcileState(ctx, instance, tidev1.PhaseRunning, tidev1.PhaseFailed); err != nil {
		return err
	}

	// update endpoints
	currentEps := make(map[string]*manifestDom.Endpoint)
	for _, ep := range inst.Manifest.Endpoints {
		currentEps[ep.Name] = ep
	}

	for _, epStatus := range instance.Status.Endpoints {
		if ep, ok := currentEps[epStatus.Name]; ok {
			ep.Hostname = epStatus.Hostname
			ep.Port = epStatus.Port
		}
	}

	return nil
}

func (is *instanceSvc) Stop(ctx context.Context, inst *instanceDom.Instance) error {
	if err := is.client.DeleteInstance(ctx, inst.Name(), utils.GetConfig().Instances.Namespace); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceDeletionFailed, "", err, nil)
	}

	// TODO wait for the instance to be deleted
	return nil
}

func (is *instanceSvc) Extend(ctx context.Context, inst *instanceDom.Instance) error {
	instance := toTideInstance(ctx, inst)

	if err := is.client.UpdateInstance(ctx, instance); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceUpdateFailed, "failed to extend instance expiry", err, nil)
	}

	return nil
}

func New(ctx context.Context) (instanceDom.Service, error) {
	k8sClient, err := k8sInfra.NewK8sClient(
		&http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	)
	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrK8sConnectionFailed, "", err, nil)
	}

	return &instanceSvc{
		client: sdk.New(k8sClient),
	}, nil
}
