package instance

import (
	"context"

	k8sInfra "github.com/TheAlpha16/isolet/oracle/infra/k8s"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/oracle/internal/domain/manifest"
	"github.com/TheAlpha16/isolet/oracle/utils"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"github.com/TheAlpha16/isolet/tide/sdk"

	"go.opentelemetry.io/otel"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var instanceTracer = otel.Tracer("api/external/instance")

type instanceSvc struct {
	client sdk.Handler
}

func (is *instanceSvc) Start(ctx context.Context, inst *instanceDom.Instance, manifest *manifestDom.Manifest) error {
	ctx, span := instanceTracer.Start(ctx, "instance.Start")
	defer span.End()

	tideInstance := toTideInstance(ctx, inst, manifest)

	if err := is.client.CreateInstance(ctx, tideInstance); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceCreationFailed, "", err, nil)
	}

	// wait for instance to reach terminal state
	if err := is.ensureInstanceReady(
		ctx, tideInstance,
		tidev1.PhaseRunning, tidev1.PhaseStaged, tidev1.PhaseFailed, tidev1.PhaseExpired, tidev1.PhaseTerminated,
	); err != nil {
		return err
	}

	inst.Endpoints = make([]*instanceDom.Endpoint, 0, len(tideInstance.Status.Endpoints))
	for _, epStatus := range tideInstance.Status.Endpoints {
		ep := fromTideEndpointStatus(epStatus)
		ep.TeamID = inst.TeamID
		inst.Endpoints = append(inst.Endpoints, ep)
	}

	return nil
}

func (is *instanceSvc) Stop(ctx context.Context, inst *instanceDom.Instance) error {
	ctx, span := instanceTracer.Start(ctx, "instance.Stop")
	defer span.End()

	namespace := utils.GetConfig().Instances.Namespace

	if err := is.client.DeleteInstance(ctx, inst.Name(), namespace); err != nil {
		if apierrors.IsNotFound(err) {
			return errorDom.Raise(ctx, errorDom.ErrK8sInstanceNotFound, "", err, nil)
		}
		return errorDom.Raise(ctx, errorDom.ErrInstanceDeletionFailed, "", err, nil)
	}

	if err := is.ensureInstanceDeleted(ctx, inst.Name(), namespace); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceDeletionFailed, "instance deletion incomplete", err, nil)
	}

	return nil
}

func (is *instanceSvc) Extend(ctx context.Context, inst *instanceDom.Instance) error {
	ctx, span := instanceTracer.Start(ctx, "instance.Extend")
	defer span.End()

	tideInstance, err := is.client.GetInstance(ctx, inst.Name(), utils.GetConfig().Instances.Namespace)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return errorDom.Raise(ctx, errorDom.ErrK8sInstanceNotFound, "", err, nil)
		}
		return errorDom.Raise(ctx, errorDom.ErrInstanceUpdateFailed, "failed to get instance for extension", err, nil)
	}

	if inst.Lifecycle.ExpiresAt != nil {
		tideInstance.Spec.Lifecycle.ExpiresAt = &metav1.Time{Time: *inst.Lifecycle.ExpiresAt}
	}
	tideInstance.Spec.Lifecycle.AllowExtension = inst.Lifecycle.AllowExtension

	if err := is.client.UpdateInstance(ctx, tideInstance); err != nil {
		if apierrors.IsNotFound(err) {
			return errorDom.Raise(ctx, errorDom.ErrK8sInstanceNotFound, "", err, nil)
		}
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
