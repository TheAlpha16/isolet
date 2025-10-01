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

package controller

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	challengesv1 "github.com/TheAlpha16/isolet/tide/api/v1"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=challenges.isolet.dev,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=challenges.isolet.dev,resources=instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=challenges.isolet.dev,resources=instances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Instance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	timeNow := v1.Now()

	// fetch the Instance resource
	var instance challengesv1.Instance
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		// resource not found, might have been deleted - nothing to do
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// error reading the object - requeue the request
		log.Error(err, "unable to fetch Instance", "instance", req.NamespacedName, "error", err)
		return ctrl.Result{}, err
	}

	// update phase of the instance if not set (new instance)
	if instance.Status.Phase == "" {
		instance.Status.Phase = challengesv1.PhasePending
		if err := r.Status().Update(ctx, &instance); err != nil {
			log.Error(err, "unable to update Instance status", "instance", req.NamespacedName, "error", err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// TODO evaluate if k8s automatically deletes the insatnce here or we need to manually call delete
		// if err := r.Delete(ctx, &instance); err != nil {
		// 	log.Error(err, "unable to delete Instance", "instance", req.NamespacedName, "error", err)
		// 	return ctrl.Result{}, err
		// }
		return ctrl.Result{}, nil
	}

	if instance.Spec.Lifecycle != nil && instance.Spec.Lifecycle.ExpiresAt != nil {
		// delete the instance if expired
		if instance.Spec.Lifecycle.ExpiresAt.Before(&timeNow) {
			// instance is expired, delete it
			if err := r.Client.Delete(ctx, &instance); err != nil {
				log.Error(err, "unable to delete expired Instance", "instance", req.NamespacedName, "error", err)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		// requeue the request to check for expiry
		return ctrl.Result{
			RequeueAfter: instance.Spec.Lifecycle.ExpiresAt.Sub(timeNow.Time),
		}, nil
	}

	if instance.Status.Phase != challengesv1.PhaseRunning {
		instance.Status.Phase = challengesv1.PhaseRunning
		if err := r.Status().Update(ctx, &instance); err != nil {
			log.Error(err, "unable to update Instance status", "instance", req.NamespacedName, "error", err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	// handle deletion if deletionTimestamp is set
	// TODO evaluate if we need to manually delete child resources, or if we can rely on owner references
	// if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
	// 	// TODO handle finalizers and cleanup of cluster resources

	// 	// remove finalizer if present
	// 	if containsString(instance.ObjectMeta.Finalizers, "instance.finalizers.challenges.isolet.dev") {
	// 		instance.ObjectMeta.Finalizers = removeString(instance.ObjectMeta.Finalizers, "instance.finalizers.challenges.isolet.dev")
	// 		if err := r.Update(ctx, &instance); err != nil {
	// 			log.Error(err, "unable to remove finalizer from Instance", "instance", req.NamespacedName, "error", err)
	// 			return ctrl.Result{}, err
	// 		}
	// 	}
	// 	return ctrl.Result{}, nil
	// }

	// TODO handle deletion if deletionTimestamp is set
	// handle all finalizers and their removal

	// TODO insert finalizers if missing

	// TODO reconsile cluster resources
	// TODO check if Deployment exists, create/update if needed
	// TODO check if Service exists, create/update if needed
	// TODO check if IngressRoute exists, create/update if needed
	// TODO resolve Hostnames for endpoints, update instance.Status.Endpoints
	// TODO update Phase to Staged / Running / Failed based on Pod readiness

	// TODO delete expired instances

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&challengesv1.Instance{}).
		Named("instance").
		Complete(r)
}
