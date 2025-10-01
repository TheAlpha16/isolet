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

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Instance.
func (v *InstanceCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	instance, ok := newObj.(*challengesv1.Instance)
	if !ok {
		return nil, fmt.Errorf("expected a Instance object for the newObj but got %T", newObj)
	}
	instancelog.Info("Validation for Instance upon update", "name", instance.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Instance.
func (v *InstanceCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	instance, ok := obj.(*challengesv1.Instance)
	if !ok {
		return nil, fmt.Errorf("expected a Instance object but got %T", obj)
	}
	instancelog.Info("Validation for Instance upon deletion", "name", instance.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
