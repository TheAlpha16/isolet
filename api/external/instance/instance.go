package instance

import (
	"context"

	k8sInfra "github.com/TheAlpha16/isolet/api/infra/k8s"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
	"github.com/TheAlpha16/isolet/api/utils"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"github.com/TheAlpha16/isolet/tide/sdk"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type instanceSvc struct {
	client sdk.Handler
}

func (is *instanceSvc) Start(ctx context.Context, inst *instanceDom.Instance, manifest *manifestDom.Manifest) error {
	tideInstance := toTideInstance(ctx, inst, manifest)

	if err := is.client.CreateInstance(ctx, tideInstance); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceCreationFailed, "", err, nil)
	}

	// wait for instance to reach terminal state
	if err := is.ensureInstanceReady(ctx, tideInstance, tidev1.PhaseRunning, tidev1.PhaseStaged, tidev1.PhaseFailed); err != nil {
		return err
	}

	inst.Endpoints = make([]*instanceDom.Endpoint, 0, len(tideInstance.Status.Endpoints))
	for _, epStatus := range tideInstance.Status.Endpoints {
		inst.Endpoints = append(inst.Endpoints, fromTideEndpointStatus(epStatus))
	}

	return nil
}

func (is *instanceSvc) Stop(ctx context.Context, inst *instanceDom.Instance) error {
	namespace := utils.GetConfig().Instances.Namespace

	if err := is.client.DeleteInstance(ctx, inst.Name(), namespace); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceDeletionFailed, "", err, nil)
	}

	if err := is.ensureInstanceDeleted(ctx, inst.Name(), namespace); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceDeletionFailed, "instance deletion incomplete", err, nil)
	}

	return nil
}

func (is *instanceSvc) Extend(ctx context.Context, inst *instanceDom.Instance) error {
	tideInstance, err := is.client.GetInstance(ctx, inst.Name(), utils.GetConfig().Instances.Namespace)
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceUpdateFailed, "failed to get instance for extension", err, nil)
	}

	if inst.Lifecycle.ExpiresAt != nil {
		tideInstance.Spec.Lifecycle.ExpiresAt = &metav1.Time{Time: *inst.Lifecycle.ExpiresAt}
	}
	tideInstance.Spec.Lifecycle.AllowExtension = inst.Lifecycle.AllowExtension

	if err := is.client.UpdateInstance(ctx, tideInstance); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceUpdateFailed, "failed to extend instance expiry", err, nil)
	}

	return nil
}

func New(ctx context.Context) (instanceDom.Service, error) {
	k8sClient, err := k8sInfra.NewK8sClient()
	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrK8sConnectionFailed, "", err, nil)
	}

	return &instanceSvc{
		client: sdk.New(k8sClient),
	}, nil
}
