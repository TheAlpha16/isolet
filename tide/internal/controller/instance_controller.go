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
	"k8s.io/apimachinery/pkg/api/meta"
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
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.1/pkg/reconcile
func (r *InstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.V(1).Info("Reconciling Instance", "namespace", req.Namespace, "name", req.Name)

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
			r.Recorder.Event(&instance, corev1.EventTypeWarning, EventReasonStatusUpdateFailed, fmt.Sprintf("Failed to set initial status: %v", err))
			return ctrl.Result{}, err
		}
		r.Recorder.Event(&instance, corev1.EventTypeNormal, EventReasonCreated, fmt.Sprintf("Instance created for challenge '%s'", instance.Spec.Challenge.Slug))
		return ctrl.Result{}, nil
	}

	// handle deletion
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("Instance is being deleted", "namespace", instance.Namespace, "name", instance.Name)
		r.Recorder.Event(&instance, corev1.EventTypeNormal, EventReasonDeleting, "Instance deletion in progress")
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
			r.Recorder.Event(&instance, corev1.EventTypeWarning, EventReasonExpired, fmt.Sprintf("Instance expired at %s", instance.Spec.Lifecycle.ExpiresAt.Time))
			if err := r.Client.Delete(ctx, &instance); err != nil {
				log.Error(err, "Failed to delete expired Instance", "namespace", instance.Namespace, "name", instance.Name)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	// capture the original phase to detect transitions later
	originalPhase := instance.Status.Phase

	// ensure child objects are in desired state
	statusChanged := false

	// 1. Reconcile Deployment
	deploymentReady, err := r.reconcileDeployment(ctx, &instance)
	if err != nil {
		log.Error(err, "failed to reconcile Deployment for Instance", "instance", req.NamespacedName)
		meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
			Type:    ConditionDeploymentReady,
			Status:  metav1.ConditionFalse,
			Reason:  "ReconciliationFailed",
			Message: err.Error(),
		})
		statusChanged = true
	} else {
		if deploymentReady {
			if meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
				Type:    ConditionDeploymentReady,
				Status:  metav1.ConditionTrue,
				Reason:  "DeploymentAvailable",
				Message: "Deployment is ready and available",
			}) {
				statusChanged = true
			}
		} else {
			// Deployment is not ready, check for Pod failures
			failed, failureReason, message := r.checkPodFailures(ctx, &instance)
			if failed {
				if meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
					Type:    ConditionDeploymentReady,
					Status:  metav1.ConditionFalse,
					Reason:  failureReason,
					Message: message,
				}) {
					statusChanged = true
				}
				// Force phase update to Failed
				instance.Status.Phase = challengesv1.PhaseFailed
				statusChanged = true // Ensure we trigger status update
			} else {
				if meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
					Type:    ConditionDeploymentReady,
					Status:  metav1.ConditionFalse,
					Reason:  "DeploymentNotReady",
					Message: "Deployment is not yet ready",
				}) {
					statusChanged = true
				}
			}
		}
	}

	// 2. Reconcile Service
	serviceReady, err := r.reconcileService(ctx, &instance)
	if err != nil {
		log.Error(err, "failed to reconcile Service for Instance", "instance", req.NamespacedName)
		meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
			Type:    ConditionServiceReady,
			Status:  metav1.ConditionFalse,
			Reason:  "ReconciliationFailed",
			Message: err.Error(),
		})
		statusChanged = true
	} else {
		if serviceReady {
			if meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
				Type:    ConditionServiceReady,
				Status:  metav1.ConditionTrue,
				Reason:  "ServiceAvailable",
				Message: "Service is ready and available",
			}) {
				statusChanged = true
			}
		} else {
			if meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
				Type:    ConditionServiceReady,
				Status:  metav1.ConditionFalse,
				Reason:  "ServiceNotReady",
				Message: "Service is not yet ready or not needed",
			}) {
				statusChanged = true
			}
		}
	}

	// 3. Reconcile Ingress
	ingressReady, err := r.reconcileIngressRoute(ctx, &instance)
	if err != nil {
		log.Error(err, "failed to reconcile Ingress for Instance", "instance", req.NamespacedName)
		meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
			Type:    ConditionIngressReady,
			Status:  metav1.ConditionFalse,
			Reason:  "ReconciliationFailed",
			Message: err.Error(),
		})
		statusChanged = true
	} else {
		if ingressReady {
			if meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
				Type:    ConditionIngressReady,
				Status:  metav1.ConditionTrue,
				Reason:  "IngressAvailable",
				Message: "Ingress is ready and available",
			}) {
				statusChanged = true
			}
		} else {
			if meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
				Type:    ConditionIngressReady,
				Status:  metav1.ConditionFalse,
				Reason:  "IngressNotReady",
				Message: "Ingress is not yet ready or not needed",
			}) {
				statusChanged = true
			}
		}
	}

	// Update Instance phase based on child resource status
	newPhase := instance.Status.Phase
	if instance.Status.Phase != challengesv1.PhaseFailed {
		newPhase = r.determinePhase(&instance, deploymentReady, serviceReady, ingressReady)
	} else {
		if deploymentReady && serviceReady && ingressReady {
			newPhase = challengesv1.PhaseRunning
		}
	}

	if originalPhase != newPhase {
		instance.Status.Phase = newPhase
		statusChanged = true
		log.Info("Instance phase transition",
			"instance", instance.Name,
			"oldPhase", originalPhase,
			"newPhase", newPhase,
			"deploymentReady", deploymentReady,
			"serviceReady", serviceReady,
			"ingressReady", ingressReady)

		// Emit events for significant phase changes
		switch newPhase {
		case challengesv1.PhaseStaged:
			r.Recorder.Event(&instance, corev1.EventTypeNormal, EventReasonInstanceStaged, "Instance staged and ready to be activated")
		case challengesv1.PhaseRunning:
			r.Recorder.Event(&instance, corev1.EventTypeNormal, EventReasonInstanceRunning, "Instance is now running and fully accessible")
		case challengesv1.PhaseFailed:
			r.Recorder.Event(&instance, corev1.EventTypeWarning, EventReasonInstanceFailed,
				fmt.Sprintf("Instance has failed (deployment: %v, service: %v, ingress: %v)",
					deploymentReady, serviceReady, ingressReady))
		}
	}

	// Update status if anything changed
	if statusChanged {
		latestInstance := &challengesv1.Instance{}
		if err := r.Get(ctx, req.NamespacedName, latestInstance); err != nil {
			log.Error(err, "Unable to refetch Instance before status patch", "instance", req.NamespacedName)
			return ctrl.Result{}, err
		}

		// Apply our status changes to the latest version
		patch := client.MergeFrom(latestInstance.DeepCopy())
		latestInstance.Status.Phase = instance.Status.Phase
		latestInstance.Status.Conditions = instance.Status.Conditions
		latestInstance.Status.Endpoints = instance.Status.Endpoints

		if err := r.Status().Patch(ctx, latestInstance, patch); err != nil {
			log.Error(err, "Unable to patch Instance status", "instance", req.NamespacedName)
			// Don't return error - let it requeue naturally and retry
			return ctrl.Result{RequeueAfter: time.Second * 2}, nil
		}
		log.Info("Successfully updated Instance status",
			"instance", instance.Name,
			"phase", latestInstance.Status.Phase,
			"endpoints", len(latestInstance.Status.Endpoints))
	}

	// Calculate requeue interval based on the nearest important event
	var requeueAfter time.Duration

	// 1. Poll for Pod failures if Pending (Deployment not ready)
	if instance.Status.Phase == challengesv1.PhasePending {
		requeueAfter = 10 * time.Second
	}

	// 2. Wait for Staged -> Running transition (AvailableAt)
	if instance.Status.Phase == challengesv1.PhaseStaged &&
		instance.Spec.Lifecycle != nil &&
		instance.Spec.Lifecycle.AvailableAt != nil {

		availableIn := instance.Spec.Lifecycle.AvailableAt.Sub(timeNow.Time)
		if availableIn > 0 {
			if requeueAfter == 0 || availableIn < requeueAfter {
				requeueAfter = availableIn
			}
		} else {
			requeueAfter = 1 * time.Second
		}
	}

	// 3. Wait for Expiry (ExpiresAt)
	if instance.Spec.Lifecycle != nil && instance.Spec.Lifecycle.ExpiresAt != nil {
		expiresIn := instance.Spec.Lifecycle.ExpiresAt.Sub(timeNow.Time)
		if expiresIn > 0 {
			if requeueAfter == 0 || expiresIn < requeueAfter {
				requeueAfter = expiresIn
			}
		}
	}

	if requeueAfter > 0 {
		log.V(1).Info("Requeuing Instance", "instance", instance.Name, "after", requeueAfter, "phase", instance.Status.Phase)
		return ctrl.Result{RequeueAfter: requeueAfter}, nil
	}

	return ctrl.Result{}, nil
}

// reconcileDeployment ensures the Deployment for the Instance exists and is up-to-date.
// Returns (ready bool, error) where ready indicates if the Deployment is available.
func (r *InstanceReconciler) reconcileDeployment(ctx context.Context, instance *challengesv1.Instance) (bool, error) {
	log := logf.FromContext(ctx)

	// Define the desired Deployment
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("deployment-%s", instance.Name),
			Namespace: instance.Namespace,
			Labels: map[string]string{
				// standard labels
				LabelAppName:      instance.Name,
				LabelAppPartOf:    "instance",
				LabelAppManagedBy: "tide-controller",
				LabelAppComponent: "deployment",

				// tide specific labels
				LabelChallengeID:   instance.Name,
				LabelChallengeSlug: instance.Spec.Challenge.Slug,
				LabelChallengeCID:  strconv.FormatInt(instance.Spec.Challenge.ID, 10),
				LabelChallengeType: string(instance.Spec.Challenge.Type),
				LabelTeamID: func() string {
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
					LabelAppName:      instance.Name,
					LabelAppComponent: "deployment",
					LabelChallengeID:  instance.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						LabelAppName:      instance.Name,
						LabelAppComponent: "deployment",
						LabelChallengeID:  instance.Name,
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

	// Apply the Deployment (Server Side Apply)
	if err := r.Patch(ctx, deployment, client.Apply, client.FieldOwner(FieldOwner), client.ForceOwnership); err != nil {
		log.Error(err, "Failed to apply Deployment", "deployment", deployment.Name)
		r.Recorder.Event(instance, corev1.EventTypeWarning, EventReasonDeploymentFailed,
			fmt.Sprintf("Failed to apply Deployment %s: %v", deployment.Name, err))
		return false, err
	}

	ready := r.isDeploymentReady(deployment)
	if ready {

	} else {
		log.V(1).Info("Deployment not yet ready",
			"deployment", deployment.Name,
			"replicas", deployment.Status.Replicas,
			"readyReplicas", deployment.Status.ReadyReplicas)
	}
	return ready, nil
}

// reconcileService ensures the Service for the Instance exists and is up-to-date.
// Returns (ready bool, error) where ready indicates if the Service is available.
func (r *InstanceReconciler) reconcileService(ctx context.Context, instance *challengesv1.Instance) (bool, error) {
	log := logf.FromContext(ctx)

	// Define the desired Service
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("svc-%s", instance.Name),
			Namespace: instance.Namespace,
			Labels: map[string]string{
				// standard labels
				LabelAppName:      instance.Name,
				LabelAppPartOf:    "instance",
				LabelAppManagedBy: "tide-controller",
				LabelAppComponent: "service",

				// tide specific labels
				LabelChallengeID:   instance.Name,
				LabelChallengeSlug: instance.Spec.Challenge.Slug,
				LabelChallengeCID:  strconv.FormatInt(instance.Spec.Challenge.ID, 10),
				LabelChallengeType: string(instance.Spec.Challenge.Type),
				LabelTeamID: func() string {
					if instance.Spec.Team != nil {
						return strconv.FormatInt(instance.Spec.Team.ID, 10)
					}
					return "dynamic"
				}(),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				LabelAppName:      instance.Name,
				LabelAppComponent: "deployment",
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

	// Handle deletion if no ports defined
	if len(service.Spec.Ports) == 0 {
		// No endpoints, ensure Service is deleted
		if err := r.Delete(ctx, service); err != nil {
			if !apierrors.IsNotFound(err) {
				log.Error(err, "Failed to delete Service", "service", service.Name)
				return false, err
			}
		}
		log.V(1).Info("Service deleted or not needed (no ports)", "service", service.Name)
		return true, nil
	}

	// Set Instance as the owner of the Service
	if err := controllerutil.SetControllerReference(instance, service, r.Scheme); err != nil {
		log.Error(err, "failed to set controller reference for Service")
		return false, err
	}

	// Apply the Service
	if err := r.Patch(ctx, service, client.Apply, client.FieldOwner(FieldOwner), client.ForceOwnership); err != nil {
		log.Error(err, "Failed to apply Service", "service", service.Name)
		r.Recorder.Event(instance, corev1.EventTypeWarning, EventReasonServiceFailed,
			fmt.Sprintf("Failed to apply Service %s: %v", service.Name, err))
		return false, err
	}

	return true, nil
}

// checkPodFailures checks if any pods for the instance are in a failed state.
// Returns (failed bool, reason string, message string)
func (r *InstanceReconciler) checkPodFailures(ctx context.Context, instance *challengesv1.Instance) (bool, string, string) {
	log := logf.FromContext(ctx)
	// List pods matching the instance labels
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels{
			LabelAppName: instance.Name,
		},
	}
	if err := r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods for failure check", "instance", instance.Name)
		return false, "", ""
	}

	for _, pod := range podList.Items {
		for _, status := range pod.Status.ContainerStatuses {
			if status.State.Waiting != nil {
				reason := status.State.Waiting.Reason
				if isFailureReason(reason) {
					return true, reason, fmt.Sprintf("Pod %s is waiting: %s", pod.Name, reason)
				}
			}
			if status.State.Terminated != nil {
				reason := status.State.Terminated.Reason
				if isFailureReason(reason) {
					return true, reason, fmt.Sprintf("Pod %s terminated: %s", pod.Name, reason)
				}
				if status.State.Terminated.ExitCode != 0 {
					return true, "ContainerError", fmt.Sprintf("Pod %s exited with code %d", pod.Name, status.State.Terminated.ExitCode)
				}
			}
		}
	}
	return false, "", ""
}

func isFailureReason(reason string) bool {
	switch reason {
	case "CrashLoopBackOff",
		"ErrImagePull",
		"ImagePullBackOff",
		"CreateContainerConfigError",
		"InvalidImageName",
		"CreateContainerError":
		return true
	}
	return false
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
			r.Recorder.Event(instance, corev1.EventTypeWarning, EventReasonIngressFailed,
				fmt.Sprintf("Failed to create IngressRoute %s: %v", ingressRouteName, err))
			return false, err
		}

		r.Recorder.Event(instance, corev1.EventTypeNormal, EventReasonIngressCreated,
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
				r.Recorder.Event(instance, corev1.EventTypeWarning, EventReasonIngressUpdateFailed,
					fmt.Sprintf("Failed to update IngressRoute %s: %v", ingressRouteName, err))
				return false, err
			}

			r.Recorder.Event(instance, corev1.EventTypeNormal, EventReasonIngressUpdated,
				fmt.Sprintf("Updated IngressRoute %s with %d route(s)", ingressRouteName, len(routes)))
		} else {

		}
	}

	// Update resolved endpoints in status if changed
	if !equalEndpointStatus(instance.Status.Endpoints, resolvedEndpoints) {
		instance.Status.Endpoints = resolvedEndpoints
		r.Recorder.Event(instance, corev1.EventTypeNormal, EventReasonEndpointsResolved,
			fmt.Sprintf("Resolved %d endpoint(s) for Instance", len(resolvedEndpoints)))
	}

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
		Owns(&traefikv1alpha1.IngressRoute{}).
		Named("instance").
		Complete(r)
}
