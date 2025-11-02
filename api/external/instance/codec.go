package instance

import (
	"context"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
	"github.com/TheAlpha16/isolet/api/utils/logger"
	"go.uber.org/zap"

	tide "github.com/TheAlpha16/isolet/tide/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func toTideInstance(ctx context.Context, inst *instanceDom.Instance) *tide.Instance {
	instance := tide.Instance{
		Spec: tide.InstanceSpec{
			Challenge: tide.Challenge{
				ID:    inst.Manifest.ChallengeID,
				Slug:  inst.Manifest.Slug,
				Flag:  inst.Manifest.Flag,
				Type:  tide.ChallengeType(inst.Manifest.Type),
				Image: inst.Manifest.Image,
			},
		},
	}

	if inst.TeamID != nil {
		instance.Spec.Team = &tide.Team{
			ID: *inst.TeamID,
		}
	}

	if inst.Manifest.Requests != nil {
		instance.Spec.Requests = *toCoreResource(ctx, inst.Manifest.Requests)
	}

	if inst.Manifest.Limits != nil {
		instance.Spec.Limits = *toCoreResource(ctx, inst.Manifest.Limits)
	}

	var endpoints []tide.EndpointSpec
	for _, ep := range inst.Manifest.Endpoints {
		endpoints = append(endpoints, *toTideEndpoint(ctx, ep))
	}

	lifecycle := tide.Lifecycle{
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

func toTideEndpoint(ctx context.Context, endpoint *manifestDom.Endpoint) *tide.EndpointSpec {
	return &tide.EndpointSpec{
		Name:       endpoint.Name,
		Protocol:   tide.Protocol(endpoint.Protocol),
		TargetPort: endpoint.TargetPort,
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
