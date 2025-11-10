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
	"time"

	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	challengesv1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"github.com/TheAlpha16/isolet/tide/utils"
)

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=challenges.isolet.dev,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=challenges.isolet.dev,resources=instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=challenges.isolet.dev,resources=instances/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=traefik.io,resources=ingressroutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling Instance", "namespace", req.Namespace, "name", req.Name)

	timeNow := metav1.Now()

	// fetch the Instance resource
	var instance challengesv1.Instance
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		// resource not found, might have been deleted - nothing to do
		if apierrors.IsNotFound(err) {
			log.Info("Instance not found, likely deleted", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{}, nil
		}
		// error reading the object - requeue the request
		log.Error(err, "Failed to fetch Instance", "namespace", req.Namespace, "name", req.Name)
		return ctrl.Result{}, err
	}

	// update phase of the instance if not set (new instance)
	if instance.Status.Phase == "" {
		log.Info("Initializing new Instance", "namespace", instance.Namespace, "name", instance.Name,
			"challenge", instance.Spec.Challenge.Slug, "team", getTeamID(&instance))
		instance.Status.Phase = challengesv1.PhasePending
		if err := r.Status().Update(ctx, &instance); err != nil {
			log.Error(err, "Failed to update Instance status to Pending", "namespace", instance.Namespace, "name", instance.Name)
			r.Recorder.Event(&instance, corev1.EventTypeWarning, "StatusUpdateFailed", fmt.Sprintf("Failed to set initial status: %v", err))
			return ctrl.Result{}, err
		}
		r.Recorder.Event(&instance, corev1.EventTypeNormal, "Created", fmt.Sprintf("Instance created for challenge '%s'", instance.Spec.Challenge.Slug))
		return ctrl.Result{}, nil
	}

	// handle deletion
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Instance is being deleted", "namespace", instance.Namespace, "name", instance.Name)
		r.Recorder.Event(&instance, corev1.EventTypeNormal, "Deleting", "Instance deletion in progress")
		return ctrl.Result{}, nil
	}

	if instance.Spec.Lifecycle != nil && instance.Spec.Lifecycle.ExpiresAt != nil {
		// delete the instance if expired
		if instance.Spec.Lifecycle.ExpiresAt.Before(&timeNow) {
			// instance is expired, delete it
			log.Info("Instance has expired, deleting",
				"namespace", instance.Namespace,
				"name", instance.Name,
				"expiresAt", instance.Spec.Lifecycle.ExpiresAt.Time)
			r.Recorder.Event(&instance, corev1.EventTypeWarning, "Expired", fmt.Sprintf("Instance expired at %s", instance.Spec.Lifecycle.ExpiresAt.Time))
			if err := r.Client.Delete(ctx, &instance); err != nil {
				log.Error(err, "Failed to delete expired Instance", "namespace", instance.Namespace, "name", instance.Name)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	// ensure child objects are in desired state
	statusChanged := false

	// 1. Reconcile Deployment
	deploymentReady, err := r.reconcileDeployment(ctx, &instance)
	if err != nil {
		log.Error(err, "failed to reconcile Deployment for Instance", "instance", req.NamespacedName)
		if r.setCondition(&instance, "DeploymentReady", metav1.ConditionFalse, "ReconciliationFailed", err.Error()) {
			statusChanged = true
		}
	} else {
		if deploymentReady {
			if r.setCondition(&instance, "DeploymentReady", metav1.ConditionTrue, "DeploymentAvailable", "Deployment is ready and available") {
				statusChanged = true
			}
		} else {
			if r.setCondition(&instance, "DeploymentReady", metav1.ConditionFalse, "DeploymentNotReady", "Deployment is not yet ready") {
				statusChanged = true
			}
		}
	}

	// 2. Reconcile Service
	serviceReady, err := r.reconcileService(ctx, &instance)
	if err != nil {
		log.Error(err, "failed to reconcile Service for Instance", "instance", req.NamespacedName)
		if r.setCondition(&instance, "ServiceReady", metav1.ConditionFalse, "ReconciliationFailed", err.Error()) {
			statusChanged = true
		}
	} else {
		if serviceReady {
			if r.setCondition(&instance, "ServiceReady", metav1.ConditionTrue, "ServiceAvailable", "Service is ready and available") {
				statusChanged = true
			}
		} else {
			if r.setCondition(&instance, "ServiceReady", metav1.ConditionFalse, "ServiceNotReady", "Service is not yet ready or not needed") {
				statusChanged = true
			}
		}
	}

	// 3. Reconcile Ingress
	ingressReady, err := r.reconcileIngressRoute(ctx, &instance)
	if err != nil {
		log.Error(err, "failed to reconcile Ingress for Instance", "instance", req.NamespacedName)
		if r.setCondition(&instance, "IngressReady", metav1.ConditionFalse, "ReconciliationFailed", err.Error()) {
			statusChanged = true
		}
	} else {
		if ingressReady {
			if r.setCondition(&instance, "IngressReady", metav1.ConditionTrue, "IngressAvailable", "Ingress is ready and available") {
				statusChanged = true
			}
		} else {
			if r.setCondition(&instance, "IngressReady", metav1.ConditionFalse, "IngressNotReady", "Ingress is not yet ready or not needed") {
				statusChanged = true
			}
		}
	}

	// Update Instance phase based on child resource status
	newPhase := r.determinePhase(&instance, deploymentReady, serviceReady, ingressReady)
	if instance.Status.Phase != newPhase {
		oldPhase := instance.Status.Phase
		instance.Status.Phase = newPhase
		statusChanged = true
		log.Info("Instance phase transition",
			"instance", instance.Name,
			"oldPhase", oldPhase,
			"newPhase", newPhase,
			"deploymentReady", deploymentReady,
			"serviceReady", serviceReady,
			"ingressReady", ingressReady)

		// Emit events for significant phase changes
		switch newPhase {
		case challengesv1.PhaseStaged:
			r.Recorder.Event(&instance, corev1.EventTypeNormal, "InstanceStaged", "Instance staged and ready to be activated")
		case challengesv1.PhaseRunning:
			r.Recorder.Event(&instance, corev1.EventTypeNormal, "InstanceRunning", "Instance is now running and fully accessible")
		case challengesv1.PhaseFailed:
			r.Recorder.Event(&instance, corev1.EventTypeWarning, "InstanceFailed",
				fmt.Sprintf("Instance has failed (deployment: %v, service: %v, ingress: %v)",
					deploymentReady, serviceReady, ingressReady))
		}
	}

	// Update status if anything changed
	if statusChanged {
		// Refetch the Instance to get the latest resourceVersion before updating status
		// This prevents "object has been modified" conflicts
		latestInstance := &challengesv1.Instance{}
		if err := r.Get(ctx, req.NamespacedName, latestInstance); err != nil {
			log.Error(err, "Unable to refetch Instance before status patch", "instance", req.NamespacedName)
			return ctrl.Result{}, err
		}

		// Apply our status changes to the latest version
		patch := client.MergeFrom(latestInstance.DeepCopy())
		latestInstance.Status.Phase = instance.Status.Phase
		latestInstance.Status.Conditions = instance.Status.Conditions

		if err := r.Status().Patch(ctx, latestInstance, patch); err != nil {
			log.Error(err, "Unable to patch Instance status", "instance", req.NamespacedName)
			// Don't return error - let it requeue naturally and retry
			return ctrl.Result{RequeueAfter: time.Second * 2}, nil
		}
		log.Info("Successfully updated Instance status",
			"instance", instance.Name,
			"phase", latestInstance.Status.Phase)
	}

	// requeue in case expiry is set
	if instance.Spec.Lifecycle != nil && instance.Spec.Lifecycle.ExpiresAt != nil {
		return ctrl.Result{
			RequeueAfter: instance.Spec.Lifecycle.ExpiresAt.Sub(timeNow.Time),
		}, nil
	}

	// TODO check if Ingress exists, create/update if needed
	// TODO resolve Hostnames for endpoints, update instance.Status.Endpoints
	// TODO update Phase to Staged / Running / Failed based on Pod readiness

	return ctrl.Result{}, nil
}

// reconcileDeployment ensures the Deployment for the Instance exists and is up-to-date.
// Returns (ready bool, error) where ready indicates if the Deployment is available.
func (r *InstanceReconciler) reconcileDeployment(ctx context.Context, instance *challengesv1.Instance) (bool, error) {
	log := logf.FromContext(ctx)

	// Define the desired Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("deployment-%s", instance.Name),
			Namespace: instance.Namespace,
			Labels: map[string]string{
				// standard labels
				"app.kubernetes.io/name":       instance.Name,
				"app.kubernetes.io/part-of":    "instance",
				"app.kubernetes.io/managed-by": "tide-controller",
				"app.kubernetes.io/component":  "deployment",

				// tide specific labels
				"challenges.isolet.dev/id":           instance.Name,
				"challenges.isolet.dev/challenge":    instance.Spec.Challenge.Slug,
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
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":      instance.Name,
					"app.kubernetes.io/component": "deployment",
					"challenges.isolet.dev/id":    instance.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":      instance.Name,
						"app.kubernetes.io/component": "deployment",
						"challenges.isolet.dev/id":    instance.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "challenge",
							Image: instance.Spec.Challenge.Image,
							Env: func() []corev1.EnvVar {
								var envs []corev1.EnvVar
								if instance.Spec.Challenge.Flag != nil {
									envs = append(envs, corev1.EnvVar{
										Name:  "FLAG",
										Value: *instance.Spec.Challenge.Flag,
									})
								}
								return envs
							}(),
							Ports: func() []corev1.ContainerPort {
								var ports []corev1.ContainerPort
								for _, ep := range instance.Spec.Endpoints {
									ports = append(ports, corev1.ContainerPort{
										Name:          ep.Name,
										ContainerPort: ep.TargetPort,
									})
								}
								return ports
							}(),
							Resources: corev1.ResourceRequirements{
								Requests: instance.Spec.Requests,
								Limits:   instance.Spec.Limits,
							},
						},
					},
				},
			},
		},
	}

	// Set Instance as the owner of the Deployment
	if err := controllerutil.SetControllerReference(instance, deployment, r.Scheme); err != nil {
		log.Error(err, "failed to set controller reference for Deployment")
		return false, err
	}

	// Check if the Deployment already exists
	foundDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}, foundDeployment)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Deployment doesn't exist, create it
			log.Info("Creating Deployment for Instance",
				"deployment", deployment.Name,
				"instance", instance.Name,
				"image", instance.Spec.Challenge.Image,
				"challenge", instance.Spec.Challenge.Slug)
			if err := r.Create(ctx, deployment); err != nil {
				log.Error(err, "Failed to create Deployment", "deployment", deployment.Name)
				r.Recorder.Event(instance, corev1.EventTypeWarning, "DeploymentCreationFailed",
					fmt.Sprintf("Failed to create Deployment %s: %v", deployment.Name, err))
				return false, err
			}
			log.Info("Successfully created Deployment", "deployment", deployment.Name)
			r.Recorder.Event(instance, corev1.EventTypeNormal, "DeploymentCreated",
				fmt.Sprintf("Created Deployment %s with image %s", deployment.Name, instance.Spec.Challenge.Image))
			// Deployment created but not yet ready
			return false, nil
		}
		// Error reading the Deployment
		log.Error(err, "Failed to get Deployment", "deployment", deployment.Name)
		return false, err
	}

	// ensure the deployment spec is up to date
	if !equalDeploymentSpec(&foundDeployment.Spec, &deployment.Spec) {
		log.Info("Deployment spec differs, updating",
			"deployment", foundDeployment.Name,
			"currentReplicas", foundDeployment.Spec.Replicas,
			"desiredReplicas", deployment.Spec.Replicas,
			"image", instance.Spec.Challenge.Image)
		foundDeployment.Spec = deployment.Spec
		if err := r.Update(ctx, foundDeployment); err != nil {
			log.Error(err, "Failed to update Deployment", "deployment", foundDeployment.Name)
			r.Recorder.Event(instance, corev1.EventTypeWarning, "DeploymentUpdateFailed",
				fmt.Sprintf("Failed to update Deployment %s: %v", foundDeployment.Name, err))
			return false, err
		}
		log.Info("Successfully updated Deployment", "deployment", foundDeployment.Name)
		r.Recorder.Event(instance, corev1.EventTypeNormal, "DeploymentUpdated",
			fmt.Sprintf("Updated Deployment %s to match desired spec", foundDeployment.Name))
		// Deployment updated but may not be ready yet
		return false, nil
	}

	// Check readiness
	ready := r.isDeploymentReady(foundDeployment)
	if ready {
		log.V(1).Info("Deployment is ready", "deployment", foundDeployment.Name)
	} else {
		log.Info("Deployment not yet ready",
			"deployment", foundDeployment.Name,
			"replicas", foundDeployment.Status.Replicas,
			"readyReplicas", foundDeployment.Status.ReadyReplicas)
	}
	return ready, nil
}

// reconcileService ensures the Service for the Instance exists and is up-to-date.
// Returns (ready bool, error) where ready indicates if the Service is available.
func (r *InstanceReconciler) reconcileService(ctx context.Context, instance *challengesv1.Instance) (bool, error) {
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
				"challenges.isolet.dev/challenge":    instance.Spec.Challenge.Slug,
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
		return false, err
	}

	// Check if the Service already exists
	foundService := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Service doesn't exist, create it if there are endpoints defined
			if len(service.Spec.Ports) == 0 {
				log.Info("No endpoints defined for Instance, skipping Service creation",
					"instance", instance.Name)
				// No service needed, so consider it "ready"
				return true, nil
			}
			log.Info("Creating Service for Instance",
				"service", service.Name,
				"instance", instance.Name,
				"ports", len(service.Spec.Ports))
			if err := r.Create(ctx, service); err != nil {
				log.Error(err, "Failed to create Service", "service", service.Name)
				r.Recorder.Event(instance, corev1.EventTypeWarning, "ServiceCreationFailed",
					fmt.Sprintf("Failed to create Service %s: %v", service.Name, err))
				return false, err
			}
			log.Info("Successfully created Service", "service", service.Name)
			r.Recorder.Event(instance, corev1.EventTypeNormal, "ServiceCreated",
				fmt.Sprintf("Created Service %s with %d port(s)", service.Name, len(service.Spec.Ports)))
			// Service created and is ready (Services are immediately available)
			return true, nil
		}
		// Error reading the Service
		log.Error(err, "Failed to get Service", "service", service.Name)
		return false, err
	}

	// ensure the service spec is up to date
	if !equalServiceSpec(&foundService.Spec, &service.Spec) {
		// if new endpoints are empty, delete the service
		if len(service.Spec.Ports) == 0 {
			log.Info("No endpoints defined for Instance, deleting Service",
				"service", foundService.Name,
				"instance", instance.Name)
			if err := r.Delete(ctx, foundService); err != nil {
				log.Error(err, "Failed to delete Service", "service", foundService.Name)
				r.Recorder.Event(instance, corev1.EventTypeWarning, "ServiceDeletionFailed",
					fmt.Sprintf("Failed to delete Service %s: %v", foundService.Name, err))
				return false, err
			}
			log.Info("Successfully deleted Service", "service", foundService.Name)
			r.Recorder.Event(instance, corev1.EventTypeNormal, "ServiceDeleted",
				fmt.Sprintf("Deleted Service %s (no endpoints defined)", foundService.Name))
			return true, nil
		}
		foundService.Spec = service.Spec
		log.Info("Updating Service for Instance",
			"service", foundService.Name,
			"instance", instance.Name,
			"ports", len(service.Spec.Ports))
		if err := r.Update(ctx, foundService); err != nil {
			log.Error(err, "Failed to update Service", "service", foundService.Name)
			r.Recorder.Event(instance, corev1.EventTypeWarning, "ServiceUpdateFailed",
				fmt.Sprintf("Failed to update Service %s: %v", foundService.Name, err))
			return false, err
		}
		log.Info("Successfully updated Service", "service", foundService.Name)
		r.Recorder.Event(instance, corev1.EventTypeNormal, "ServiceUpdated",
			fmt.Sprintf("Updated Service %s with %d port(s)", foundService.Name, len(service.Spec.Ports)))
	}

	// Service exists and is ready
	log.V(1).Info("Service is ready", "service", foundService.Name)
	return true, nil
}

func equalDeploymentSpec(a, b *appsv1.DeploymentSpec) bool {
	// replicas (treat nil as 1, which is Kubernetes default)
	aReplicas := int32(1)
	if a.Replicas != nil {
		aReplicas = *a.Replicas
	}
	bReplicas := int32(1)
	if b.Replicas != nil {
		bReplicas = *b.Replicas
	}
	if aReplicas != bReplicas {
		return false
	}

	// selector match labels
	if len(a.Selector.MatchLabels) != len(b.Selector.MatchLabels) {
		return false
	}
	for k, v := range a.Selector.MatchLabels {
		if bv, exists := b.Selector.MatchLabels[k]; !exists || bv != v {
			return false
		}
	}

	// labels
	if len(a.Template.Labels) != len(b.Template.Labels) {
		return false
	}
	for k, v := range a.Template.Labels {
		if bv, exists := b.Template.Labels[k]; !exists || bv != v {
			return false
		}
	}

	// containers
	if len(a.Template.Spec.Containers) != len(b.Template.Spec.Containers) {
		return false
	}
	containerMap := make(map[string]corev1.Container)
	for _, c := range a.Template.Spec.Containers {
		containerMap[c.Name] = c
	}
	for _, c := range b.Template.Spec.Containers {
		// name, image
		ac, exists := containerMap[c.Name]
		if !exists || ac.Image != c.Image {
			return false
		}

		// env
		if len(ac.Env) != len(c.Env) {
			return false
		}
		envMap := make(map[string]corev1.EnvVar)
		for _, e := range ac.Env {
			envMap[e.Name] = e
		}
		for _, e := range c.Env {
			if ae, exists := envMap[e.Name]; !exists || ae.Value != e.Value {
				return false
			}
		}

		// ports
		if len(ac.Ports) != len(c.Ports) {
			return false
		}
		portMap := make(map[string]corev1.ContainerPort)
		for _, p := range ac.Ports {
			portMap[p.Name] = p
		}
		for _, p := range c.Ports {
			if ap, exists := portMap[p.Name]; !exists || ap.ContainerPort != p.ContainerPort {
				return false
			}
		}

		// resources (on container level)
		// requests
		if len(ac.Resources.Requests) != len(c.Resources.Requests) {
			return false
		}
		for k, v := range ac.Resources.Requests {
			if cv, exists := c.Resources.Requests[k]; !exists || cv.Cmp(v) != 0 {
				return false
			}
		}

		// limits
		if len(ac.Resources.Limits) != len(c.Resources.Limits) {
			return false
		}
		for k, v := range ac.Resources.Limits {
			if cv, exists := c.Resources.Limits[k]; !exists || cv.Cmp(v) != 0 {
				return false
			}
		}
	}

	return true
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

// isDeploymentReady checks if a Deployment is available and ready.
func (r *InstanceReconciler) isDeploymentReady(deployment *appsv1.Deployment) bool {
	// Get desired replica count (default to 1 if not specified)
	desiredReplicas := int32(1)
	if deployment.Spec.Replicas != nil {
		desiredReplicas = *deployment.Spec.Replicas
	}

	// Check if all desired replicas are ready and available
	if deployment.Status.ReadyReplicas >= desiredReplicas &&
		deployment.Status.AvailableReplicas >= desiredReplicas &&
		deployment.Status.UpdatedReplicas >= desiredReplicas {
		return true
	}
	return false
}

// determinePhase determines the Instance phase based on child resource status.
func (r *InstanceReconciler) determinePhase(instance *challengesv1.Instance, deploymentReady, serviceReady, ingressReady bool) challengesv1.Phase {
	// If Deployment is not ready, Instance is Pending
	if !deploymentReady {
		return challengesv1.PhasePending
	}

	// Check if we should be in Staged phase (availableAt in the future)
	if instance.Spec.Lifecycle != nil && instance.Spec.Lifecycle.AvailableAt != nil {
		if instance.Spec.Lifecycle.AvailableAt.After(metav1.Now().Time) {
			return challengesv1.PhaseStaged
		}
	}

	// If Deployment, Service and Ingress are ready (or not needed), Instance is Running
	if deploymentReady && serviceReady && ingressReady {
		return challengesv1.PhaseRunning
	}

	// Default to Pending if we can't determine
	return challengesv1.PhasePending
}

// setCondition adds or updates a condition in the Instance status.
// Returns true if the condition was changed.
func (r *InstanceReconciler) setCondition(instance *challengesv1.Instance, conditionType string, status metav1.ConditionStatus, reason, message string) bool {
	now := metav1.Now()

	// Find existing condition
	for i, condition := range instance.Status.Conditions {
		if condition.Type == conditionType {
			// Update existing condition only if it changed
			if condition.Status != status || condition.Reason != reason || condition.Message != message {
				instance.Status.Conditions[i].Status = status
				instance.Status.Conditions[i].Reason = reason
				instance.Status.Conditions[i].Message = message
				instance.Status.Conditions[i].LastTransitionTime = now
				instance.Status.Conditions[i].ObservedGeneration = instance.Generation
				return true // Condition changed
			}
			return false // Condition unchanged
		}
	}

	// Add new condition
	instance.Status.Conditions = append(instance.Status.Conditions, metav1.Condition{
		Type:               conditionType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: now,
		ObservedGeneration: instance.Generation,
	})
	return true // New condition added
}

// reconcileIngressRoute ensures IngressRoute resources exist for HTTP/HTTPS endpoints.
// Returns (ready bool, err error) where ready indicates if IngressRoutes are properly configured.
func (r *InstanceReconciler) reconcileIngressRoute(ctx context.Context, instance *challengesv1.Instance) (bool, error) {
	log := logf.FromContext(ctx)

	// Filter endpoints that need IngressRoute (only HTTP/HTTPS)
	var httpEndpoints []challengesv1.EndpointSpec
	for _, endpoint := range instance.Spec.Endpoints {
		if endpoint.Protocol == challengesv1.ProtocolHTTP || endpoint.Protocol == challengesv1.ProtocolHTTPS {
			httpEndpoints = append(httpEndpoints, endpoint)
		}
	}

	// If no HTTP/HTTPS endpoints, ensure no IngressRoute exists and return ready
	if len(httpEndpoints) == 0 {
		log.Info("No HTTP/HTTPS endpoints defined for Instance, skipping IngressRoute creation", "instance", instance.Name)
		return true, nil
	}

	// For each HTTP/HTTPS endpoint, create/update an IngressRoute
	ingressRouteName := fmt.Sprintf("ir-%s", instance.Name)
	serviceName := fmt.Sprintf("svc-%s", instance.Name)
	resolvedEndpoints := []challengesv1.EndpointStatus{}

	// Define the desired IngressRoute using Traefik SDK
	ingressRoute := &traefikv1alpha1.IngressRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressRouteName,
			Namespace: instance.Namespace,
			Labels: map[string]string{
				// standard labels
				"app.kubernetes.io/name":       instance.Name,
				"app.kubernetes.io/part-of":    "instance",
				"app.kubernetes.io/managed-by": "tide-controller",
				"app.kubernetes.io/component":  "ingress",

				// tide specific labels
				"challenges.isolet.dev/id":           instance.Name,
				"challenges.isolet.dev/challenge":    instance.Spec.Challenge.Slug,
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
		Spec: traefikv1alpha1.IngressRouteSpec{
			EntryPoints: []string{"web", "websecure"},
			TLS: &traefikv1alpha1.TLS{
				SecretName: "challenge-certs",
			},
		},
	}

	// Traefik routes
	var routes []traefikv1alpha1.Route
	if len(httpEndpoints) == 1 {
		// Single endpoint, route directly
		resEndpoint := challengesv1.EndpointStatus{
			EndpointSpec: httpEndpoints[0],
			Hostname:     utils.Ptr(fmt.Sprintf("%s.%s.isolet.dev", instance.Name, instance.Spec.Challenge.Slug)), // using the match rule as hostname
			Ready:        true,
		}
		resolvedEndpoints = append(resolvedEndpoints, resEndpoint)
		routes = []traefikv1alpha1.Route{
			{
				Kind:  "Rule",
				Match: fmt.Sprintf("Host(`%s`)", *resEndpoint.Hostname),
				Services: []traefikv1alpha1.Service{
					{LoadBalancerSpec: traefikv1alpha1.LoadBalancerSpec{
						Name: serviceName,
						Port: intstr.FromInt32(httpEndpoints[0].TargetPort),
					}},
				},
			},
		}
	} else {
		// Multiple endpoints, include endpoint name in path
		for _, ep := range httpEndpoints {
			resEP := challengesv1.EndpointStatus{
				EndpointSpec: ep,
				Hostname:     utils.Ptr(fmt.Sprintf("%s-%s.%s.isolet.dev", ep.Name, instance.Name, instance.Spec.Challenge.Slug)),
				Ready:        true,
			}
			resolvedEndpoints = append(resolvedEndpoints, resEP)
			route := traefikv1alpha1.Route{
				Kind:  "Rule",
				Match: fmt.Sprintf("Host(`%s-%s.%s.isolet.dev`)", ep.Name, instance.Name, instance.Spec.Challenge.Slug),
				Services: []traefikv1alpha1.Service{
					{
						LoadBalancerSpec: traefikv1alpha1.LoadBalancerSpec{
							Name: serviceName,
							Port: intstr.FromInt32(ep.TargetPort),
						},
					},
				},
			}
			routes = append(routes, route)
		}
	}
	ingressRoute.Spec.Routes = routes

	// Set the Instance as the owner
	if err := controllerutil.SetControllerReference(instance, ingressRoute, r.Scheme); err != nil {
		log.Error(err, "failed to set controller reference on IngressRoute", "ingressRoute", ingressRouteName)
		return false, err
	}

	// Check if IngressRoute already exists
	foundIngressRoute := &traefikv1alpha1.IngressRoute{}
	err := r.Get(ctx, client.ObjectKey{Name: ingressRouteName, Namespace: instance.Namespace}, foundIngressRoute)
	if err != nil && apierrors.IsNotFound(err) {
		// Create the IngressRoute
		log.Info("Creating IngressRoute for Instance",
			"ingressRoute", ingressRouteName,
			"instance", instance.Name,
			"endpoints", len(httpEndpoints),
			"routes", len(routes))
		if err := r.Create(ctx, ingressRoute); err != nil {
			log.Error(err, "Failed to create IngressRoute", "ingressRoute", ingressRouteName)
			r.Recorder.Event(instance, corev1.EventTypeWarning, "IngressRouteCreationFailed",
				fmt.Sprintf("Failed to create IngressRoute %s: %v", ingressRouteName, err))
			return false, err
		}
		log.Info("Successfully created IngressRoute", "ingressRoute", ingressRouteName)
		r.Recorder.Event(instance, corev1.EventTypeNormal, "IngressRouteCreated",
			fmt.Sprintf("Created IngressRoute %s with %d route(s)", ingressRouteName, len(routes)))
	} else if err != nil {
		log.Error(err, "Failed to get IngressRoute", "ingressRoute", ingressRouteName)
		return false, err
	} else {
		// IngressRoute exists, check if update is needed
		if !equalIngressRouteSpec(&foundIngressRoute.Spec, &ingressRoute.Spec) {
			log.Info("Updating IngressRoute for Instance",
				"ingressRoute", ingressRouteName,
				"instance", instance.Name,
				"routes", len(routes))
			foundIngressRoute.Spec = ingressRoute.Spec
			foundIngressRoute.Labels = ingressRoute.Labels
			if err := r.Update(ctx, foundIngressRoute); err != nil {
				log.Error(err, "Failed to update IngressRoute", "ingressRoute", ingressRouteName)
				r.Recorder.Event(instance, corev1.EventTypeWarning, "IngressRouteUpdateFailed",
					fmt.Sprintf("Failed to update IngressRoute %s: %v", ingressRouteName, err))
				return false, err
			}
			log.Info("Successfully updated IngressRoute", "ingressRoute", ingressRouteName)
			r.Recorder.Event(instance, corev1.EventTypeNormal, "IngressRouteUpdated",
				fmt.Sprintf("Updated IngressRoute %s with %d route(s)", ingressRouteName, len(routes)))
		} else {
			log.V(1).Info("IngressRoute is up to date", "ingressRoute", ingressRouteName)
		}
	}

	// Update resolved endpoints in status if changed
	if !equalEndpointStatus(instance.Status.Endpoints, resolvedEndpoints) {
		log.Info("Updating Instance status with resolved endpoints",
			"instance", instance.Name,
			"endpoints", len(resolvedEndpoints))
		instance.Status.Endpoints = resolvedEndpoints
		if err := r.Status().Update(ctx, instance); err != nil {
			log.Error(err, "Failed to update Instance status with resolved endpoints", "instance", instance.Name)
			r.Recorder.Event(instance, corev1.EventTypeWarning, "EndpointResolutionFailed",
				fmt.Sprintf("Failed to update resolved endpoints: %v", err))
			return false, err
		}
		r.Recorder.Event(instance, corev1.EventTypeNormal, "EndpointsResolved",
			fmt.Sprintf("Resolved %d endpoint(s) for Instance", len(resolvedEndpoints)))
	}

	log.V(1).Info("IngressRoute is ready", "ingressRoute", ingressRouteName)
	return true, nil
}

// equalIngressRouteSpec compares two IngressRouteSpec for equality
func equalIngressRouteSpec(a, b *traefikv1alpha1.IngressRouteSpec) bool {
	// Compare entry points
	if len(a.EntryPoints) != len(b.EntryPoints) {
		return false
	}
	for i := range a.EntryPoints {
		if a.EntryPoints[i] != b.EntryPoints[i] {
			return false
		}
	}

	// Compare routes
	if len(a.Routes) != len(b.Routes) {
		return false
	}
	for i := range a.Routes {
		if a.Routes[i].Match != b.Routes[i].Match || a.Routes[i].Kind != b.Routes[i].Kind {
			return false
		}
		if len(a.Routes[i].Services) != len(b.Routes[i].Services) {
			return false
		}
		for j := range a.Routes[i].Services {
			if a.Routes[i].Services[j].Name != b.Routes[i].Services[j].Name ||
				a.Routes[i].Services[j].Port != b.Routes[i].Services[j].Port {
				return false
			}
		}
	}

	// Compare TLS
	if (a.TLS == nil) != (b.TLS == nil) {
		return false
	}
	if a.TLS != nil && b.TLS != nil {
		if a.TLS.SecretName != b.TLS.SecretName {
			return false
		}
	}

	return true
}

func equalEndpointStatus(a, b []challengesv1.EndpointStatus) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]challengesv1.EndpointStatus)
	for _, ep := range a {
		aMap[ep.Name] = ep
	}
	for _, ep := range b {
		if aep, exists := aMap[ep.Name]; !exists ||
			aep.Protocol != ep.Protocol ||
			aep.TargetPort != ep.TargetPort ||
			aep.Hostname == nil || ep.Hostname == nil || *aep.Hostname != *ep.Hostname ||
			aep.Ready != ep.Ready {
			return false
		}
	}
	return true
}

// getTeamID returns the team ID as a string, or "dynamic" if not set
func getTeamID(instance *challengesv1.Instance) string {
	if instance.Spec.Team != nil {
		return strconv.FormatInt(instance.Spec.Team.ID, 10)
	}
	return "dynamic"
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&challengesv1.Instance{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Named("instance").
		Complete(r)
}
