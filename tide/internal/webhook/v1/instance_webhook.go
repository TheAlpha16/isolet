/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	challengesv1 "github.com/TheAlpha16/isolet/tide/api/v1"
)

// nolint:unused
// log is for logging in this package.
var instancelog = logf.Log.WithName("instance-resource")

// SetupInstanceWebhookWithManager registers the webhook for Instance in the manager.
func SetupInstanceWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&challengesv1.Instance{}).
		WithValidator(&InstanceCustomValidator{}).
		WithDefaulter(&InstanceCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-challenges-isolet-dev-v1-instance,mutating=true,failurePolicy=fail,sideEffects=None,groups=challenges.isolet.dev,resources=instances,verbs=create;update,versions=v1,name=minstance-v1.kb.io,admissionReviewVersions=v1

// InstanceCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Instance when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type InstanceCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &InstanceCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Instance.
func (d *InstanceCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	instance, ok := obj.(*challengesv1.Instance)

	if !ok {
		return fmt.Errorf("expected an Instance object but got %T", obj)
	}
	instancelog.Info("Defaulting for Instance", "name", instance.GetName())

	// make sure default lifecycle is set
	if instance.Spec.Lifecycle == nil {
		instance.Spec.Lifecycle = &challengesv1.Lifecycle{}
		instance.Spec.Lifecycle.SetDefaults()
	}

	if instance.Spec.Lifecycle != nil && instance.Spec.Lifecycle.RestartPolicy == "" {
		instance.Spec.Lifecycle.RestartPolicy = corev1.RestartPolicyNever
	}

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-challenges-isolet-dev-v1-instance,mutating=false,failurePolicy=fail,sideEffects=None,groups=challenges.isolet.dev,resources=instances,verbs=create;update,versions=v1,name=vinstance-v1.kb.io,admissionReviewVersions=v1

// InstanceCustomValidator struct is responsible for validating the Instance resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type InstanceCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &InstanceCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Instance.
func (v *InstanceCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	instance, ok := obj.(*challengesv1.Instance)
	if !ok {
		return nil, fmt.Errorf("expected a Instance object but got %T", obj)
	}
	instancelog.Info("Validation for Instance upon creation", "name", instance.GetName())

	// enforce team for on-demand challenges
	if instance.Spec.Challenge.Type == challengesv1.ChallengeTypeOnDemand && instance.Spec.Team == nil {
		return nil, fmt.Errorf("team must be set for on-demand challenges")
	}

	// validate endpoints
	if err := v.validateEndpoints(instance.Spec.Endpoints); err != nil {
		return nil, err
	}

	// validate lifecycle
	if err := v.validateLifecycle(instance.Spec.Lifecycle); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Instance.
func (v *InstanceCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	newInstance, ok := newObj.(*challengesv1.Instance)
	if !ok {
		return nil, fmt.Errorf("expected a Instance object for the newObj but got %T", newObj)
	}
	instancelog.Info("Validation for Instance upon update", "name", newInstance.GetName())

	oldInstance, ok := oldObj.(*challengesv1.Instance)
	if !ok {
		return nil, fmt.Errorf("expected a Instance object for the oldObj but got %T", oldObj)
	}

	// prevent changes to challenge
	if oldInstance.Spec.Challenge.ID != newInstance.Spec.Challenge.ID {
		return nil, fmt.Errorf("challenge.id is immutable")
	}
	if oldInstance.Spec.Challenge.Name != newInstance.Spec.Challenge.Name {
		return nil, fmt.Errorf("challenge.name is immutable")
	}
	if oldInstance.Spec.Challenge.Image != newInstance.Spec.Challenge.Image {
		return nil, fmt.Errorf("challenge.image is immutable")
	}
	if oldInstance.Spec.Challenge.Type != newInstance.Spec.Challenge.Type {
		return nil, fmt.Errorf("challenge.type is immutable")
	}

	// prevent changes to team
	if (oldInstance.Spec.Team == nil) != (newInstance.Spec.Team == nil) ||
		(oldInstance.Spec.Team != nil && newInstance.Spec.Team != nil && oldInstance.Spec.Team.ID != newInstance.Spec.Team.ID) {
		return nil, fmt.Errorf("team is immutable")
	}

	// validate endpoints
	if err := v.validateEndpoints(newInstance.Spec.Endpoints); err != nil {
		return nil, err
	}

	// validate lifecycle
	if err := v.validateLifecycle(newInstance.Spec.Lifecycle); err != nil {
		return nil, err
	}

	// enforce extention rules
	if err := v.enforceExtentionRules(oldInstance, newInstance); err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Instance.
func (v *InstanceCustomValidator) ValidateDelete(ctx context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (v *InstanceCustomValidator) validateEndpoints(endpoints []challengesv1.EndpointSpec) error {
	namesMap := make(map[string]struct{})

	for _, ep := range endpoints {
		// check for duplicate names
		if _, exists := namesMap[ep.Name]; exists {
			return fmt.Errorf("duplicate endpoint name: %s", ep.Name)
		}
		namesMap[ep.Name] = struct{}{}

		// validate port range
		if ep.TargetPort < 1 || ep.TargetPort > 65535 {
			return fmt.Errorf("endpoint port must be between 1 and 65535, got %d", ep.TargetPort)
		}

		// validate protocol
		switch ep.Protocol {
		case challengesv1.ProtocolHTTP, challengesv1.ProtocolHTTPS, challengesv1.ProtocolNC, challengesv1.ProtocolSSH:
			// valid
		default:
			return fmt.Errorf("endpoint protocol must be one of http, https, nc, ssh; got %s", ep.Protocol)
		}
	}
	return nil
}

func (v *InstanceCustomValidator) validateLifecycle(l *challengesv1.Lifecycle) error {
	if l == nil {
		return nil
	}

	// make sure expiry is in the future if set
	if l.ExpiresAt != nil && l.ExpiresAt.Time.Before(metav1.Now().Time) {
		return fmt.Errorf("lifecycle.expiresAt must be in the future")
	}

	// make sure availableAt is before expiresAt if both are set
	if l.AvailableAt != nil && l.ExpiresAt != nil && l.AvailableAt.Time.After(l.ExpiresAt.Time) {
		return fmt.Errorf("lifecycle.availableAt must be before lifecycle.expiresAt")
	}

	// validate restart policy
	switch l.RestartPolicy {
	case corev1.RestartPolicyAlways, corev1.RestartPolicyNever, corev1.RestartPolicyOnFailure:
		// valid
	default:
		return fmt.Errorf("lifecycle.restartPolicy must be one of Always, OnFailure, Never")
	}
	return nil
}

func (v *InstanceCustomValidator) enforceExtentionRules(oldInstance, newInstance *challengesv1.Instance) error {
	var oldExpiry, newExpiry *metav1.Time
	if oldInstance.Spec.Lifecycle != nil {
		oldExpiry = oldInstance.Spec.Lifecycle.ExpiresAt
	}
	if newInstance.Spec.Lifecycle != nil {
		newExpiry = newInstance.Spec.Lifecycle.ExpiresAt
	}

	// no change in expiry, nothing to do
	if (oldExpiry == nil && newExpiry == nil) || (oldExpiry != nil && newExpiry != nil && oldExpiry.Equal(newExpiry)) {
		return nil
	}

	// removing expiry is allowed
	if newExpiry == nil {
		return nil
	}

	// if extention is not allowed, reject
	if newInstance.Spec.Lifecycle == nil || !newInstance.Spec.Lifecycle.AllowExtension {
		if oldExpiry == nil || newExpiry.After(oldExpiry.Time) {
			return fmt.Errorf("lifecycle.expiresAt cannot be changed, as lifecycle.allowExtension is false")
		}
	}

	return nil
}
