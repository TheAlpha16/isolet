package instance

import (
	"context"

	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/oracle/internal/domain/manifest"
	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"
	"go.uber.org/zap"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func toTideInstance(ctx context.Context, inst *instanceDom.Instance, manifest *manifestDom.Manifest) *tidev1.Instance {
	instance := tidev1.Instance{
		TypeMeta: metav1.TypeMeta{
			Kind:       utils.GetConfig().K8s.InstanceKind,
			APIVersion: utils.GetConfig().K8s.InstanceAPIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      inst.Name(),
			Namespace: utils.GetConfig().Instances.Namespace,
		},
		Spec: tidev1.InstanceSpec{
			Challenge: tidev1.Challenge{
				ID:     manifest.ChallengeID,
				Slug:   manifest.Slug,
				Flag:   inst.Flag,
				Type:   tidev1.ChallengeType(manifest.Type),
				Image:  manifest.Image,
				Domain: inst.Domain,
			},
		},
	}

	if inst.TeamID != nil {
		instance.Spec.Team = &tidev1.Team{
			ID: *inst.TeamID,
		}
	}

	if manifest.Requests != nil {
		instance.Spec.Requests = *toCoreResource(ctx, manifest.Requests)
	}

	if manifest.Limits != nil {
		instance.Spec.Limits = *toCoreResource(ctx, manifest.Limits)
	}

	var endpoints []tidev1.EndpointSpec
	for _, ep := range manifest.EndpointSpecs {
		endpoints = append(endpoints, *toTideEndpointSpec(ep))
	}

	lifecycle := tidev1.Lifecycle{
		AllowExtension: inst.Lifecycle.AllowExtension,
	}
	if inst.Lifecycle.ExpiresAt != nil {
		lifecycle.ExpiresAt = &metav1.Time{Time: *inst.Lifecycle.ExpiresAt}
	}
	if inst.Lifecycle.AvailableAt != nil {
		lifecycle.AvailableAt = &metav1.Time{Time: *inst.Lifecycle.AvailableAt}
	}

	instance.Spec.Lifecycle = &lifecycle
	instance.Spec.Endpoints = endpoints

	return &instance
}

func toTideEndpointSpec(endpointSpec *manifestDom.EndpointSpec) *tidev1.EndpointSpec {
	return &tidev1.EndpointSpec{
		Name:       endpointSpec.Name,
		Protocol:   tidev1.Protocol(endpointSpec.Protocol),
		TargetPort: endpointSpec.TargetPort,
	}
}

func fromTideEndpointStatus(epStatus tidev1.EndpointStatus) *instanceDom.Endpoint {
	return &instanceDom.Endpoint{
		Name:       epStatus.Name,
		Protocol:   manifestDom.Protocol(epStatus.Protocol),
		TargetPort: epStatus.TargetPort,
		Hostname:   epStatus.Hostname,
		Port:       epStatus.Port,
		Ready:      epStatus.Ready,
	}
}

func toCoreResource(ctx context.Context, resources []*manifestDom.Resource) *corev1.ResourceList {
	list := corev1.ResourceList{}
	for _, res := range resources {
		quantity, err := resource.ParseQuantity(res.Value)
		if err != nil {
			parseErr := errorDom.Raise(ctx, errorDom.ErrInstanceInvalidResourceQuantity, "", err, common.ExtraData{"resource_id": res.ID, "value": res.Value})
			logger.GetAppLogger().Error(
				"invalid value for resource quantity",
				zap.Int64("resource_id", res.ID),
				zap.String("value", res.Value),
			)
			errorDom.RaiseToSentry(ctx, parseErr)
			continue
		}
		list[corev1.ResourceName(res.Name)] = quantity
	}
	return &list
}
