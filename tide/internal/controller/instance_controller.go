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
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

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
	timeNow := metav1.Now()

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
	}

	if instance.Status.Phase != challengesv1.PhaseRunning {
		instance.Status.Phase = challengesv1.PhaseRunning
		if err := r.Status().Update(ctx, &instance); err != nil {
			log.Error(err, "unable to update Instance status", "instance", req.NamespacedName, "error", err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// ensure child objects are in desired state
	if err := r.reconcileService(ctx, &instance); err != nil {
		log.Error(err, "failed to reconcile Service for Instance", "instance", req.NamespacedName)
		return ctrl.Result{}, err
	}

	// requeue in case expiry is set
	if instance.Spec.Lifecycle != nil && instance.Spec.Lifecycle.ExpiresAt != nil {
		return ctrl.Result{
			RequeueAfter: instance.Spec.Lifecycle.ExpiresAt.Sub(timeNow.Time),
		}, nil
	}

	// TODO reconsile cluster resources
	// TODO check if Deployment exists, create/update if needed
	// TODO check if Service exists, create/update if needed
	// TODO check if IngressRoute exists, create/update if needed
	// TODO resolve Hostnames for endpoints, update instance.Status.Endpoints
	// TODO update Phase to Staged / Running / Failed based on Pod readiness

	return ctrl.Result{}, nil
}

// reconcileService ensures the Service for the Instance exists and is up-to-date.
func (r *InstanceReconciler) reconcileService(ctx context.Context, instance *challengesv1.Instance) error {
	log := logf.FromContext(ctx)

	// Define the desired Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("svc-%s", instance.Name),
			Namespace: instance.Namespace,
			Labels: map[string]string{
				// standard labels
				"app.kubernetes.io/name":       instance.Name,
				"app.kubernetes.io/part-of":    "instance",
				"app.kubernetes.io/managed-by": "tide-controller",
				"app.kubernetes.io/component":  "service",

				// tide specific labels
				"challenges.isolet.dev/id":           instance.Name,
				"challenges.isolet.dev/challenge":    instance.Spec.Challenge.Name,
				"challenges.isolet.dev/challenge-id": strconv.FormatInt(instance.Spec.Challenge.ID, 10),
				"challenges.isolet.dev/type":         string(instance.Spec.Challenge.Type),
				"challenges.isolet.dev/team": func() string {
					if instance.Spec.Team != nil {
						return strconv.FormatInt(instance.Spec.Team.ID, 10)
					}
					return "dynamic"
				}(),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/name":      instance.Name,
				"app.kubernetes.io/component": "deployment",
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// add enpoints as service ports
	for _, endpoint := range instance.Spec.Endpoints {
		servicePort := corev1.ServicePort{
			Name: endpoint.Name,
			Port: endpoint.TargetPort,
		}
		service.Spec.Ports = append(service.Spec.Ports, servicePort)
	}

	// Set Instance as the owner of the Service
	if err := controllerutil.SetControllerReference(instance, service, r.Scheme); err != nil {
		log.Error(err, "failed to set controller reference for Service")
		return err
	}

	// Check if the Service already exists
	foundService := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Service doesn't exist, create it
			log.Info("Creating Service for Instance", "service", service.Name, "instance", instance.Name)
			if err := r.Create(ctx, service); err != nil {
				log.Error(err, "failed to create Service", "service", service.Name)
				return err
			}
			return nil
		}
		// Error reading the Service
		log.Error(err, "failed to get Service", "service", service.Name)
		return err
	}

	// Service exists, update it if needed
	log.Info("Service already exists for Instance", "service", foundService.Name, "instance", instance.Name)

	// ensure the service spec is up to date
	if !equalServiceSpec(&foundService.Spec, &service.Spec) {
		foundService.Spec = service.Spec
		log.Info("Updating Service for Instance", "service", foundService.Name, "instance", instance.Name)
		if err := r.Update(ctx, foundService); err != nil {
			log.Error(err, "failed to update Service", "service", foundService.Name)
			return err
		}
		log.Info("Successfully updated Service for Instance", "service", foundService.Name, "instance", instance.Name)
	}
	return nil
}

func equalServiceSpec(a, b *corev1.ServiceSpec) bool {
	if a.Type != b.Type {
		return false
	}
	if len(a.Ports) != len(b.Ports) {
		return false
	}
	portMap := make(map[string]corev1.ServicePort)
	for _, port := range a.Ports {
		portMap[port.Name] = port
	}
	for _, port := range b.Ports {
		if p, exists := portMap[port.Name]; !exists || p.Port != port.Port {
			return false
		}
	}
	if len(a.Selector) != len(b.Selector) {
		return false
	}
	for k, v := range a.Selector {
		if bv, exists := b.Selector[k]; !exists || bv != v {
			return false
		}
	}
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&challengesv1.Instance{}).
		Owns(&corev1.Service{}).
		Named("instance").
		Complete(r)
}
