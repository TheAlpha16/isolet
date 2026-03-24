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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	traefikv1alpha1 "github.com/traefik/traefik/v3/pkg/provider/kubernetes/crd/traefikio/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	challengesv1 "github.com/TheAlpha16/isolet/tide/api/v1"
)

var _ = Describe("Instance Controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	var (
		ctx        context.Context
		reconciler *InstanceReconciler
	)

	BeforeEach(func() {
		ctx = context.Background()
		reconciler = &InstanceReconciler{
			Client:   k8sClient,
			Scheme:   k8sClient.Scheme(),
			Recorder: record.NewFakeRecorder(100),
		}
	})

	// Helper function to create a basic test Instance
	createBasicInstance := func(name, namespace string) *challengesv1.Instance {
		return &challengesv1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: challengesv1.InstanceSpec{
				Challenge: challengesv1.Challenge{
					ID:     1,
					Slug:   "test-challenge",
					Type:   challengesv1.ChallengeTypeDynamic,
					Image:  "nginx:latest",
					Domain: "isolet.dev",
				},
				Endpoints: []challengesv1.EndpointSpec{
					{
						Name:       "http",
						Protocol:   challengesv1.ProtocolHTTP,
						TargetPort: 80,
					},
				},
			},
		}
	}

	// Helper function to simulate Deployment becoming ready
	simulateDeploymentReady := func(ctx context.Context, instanceName, namespace string) {
		deploymentName := "deployment-" + instanceName
		deployment := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{
				Name:      deploymentName,
				Namespace: namespace,
			}, deployment)
		}, timeout, interval).Should(Succeed())

		// Simulate Deployment becoming ready
		deployment.Status.Replicas = 1
		deployment.Status.ReadyReplicas = 1
		deployment.Status.AvailableReplicas = 1
		deployment.Status.UpdatedReplicas = 1
		deployment.Status.Conditions = []appsv1.DeploymentCondition{
			{
				Type:   appsv1.DeploymentAvailable,
				Status: corev1.ConditionTrue,
				Reason: "MinimumReplicasAvailable",
			},
		}
		Expect(k8sClient.Status().Update(ctx, deployment)).To(Succeed())
	}

	// Helper function to simulate Deployment failure
	simulateDeploymentFailed := func(ctx context.Context, instanceName, namespace string) {
		deploymentName := "deployment-" + instanceName
		deployment := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{
				Name:      deploymentName,
				Namespace: namespace,
			}, deployment)
		}, timeout, interval).Should(Succeed())

		// Simulate Deployment failing
		deployment.Status.Replicas = 1
		deployment.Status.ReadyReplicas = 0
		deployment.Status.AvailableReplicas = 0
		deployment.Status.UnavailableReplicas = 1
		deployment.Status.Conditions = []appsv1.DeploymentCondition{
			{
				Type:   appsv1.DeploymentProgressing,
				Status: corev1.ConditionFalse,
				Reason: "ProgressDeadlineExceeded",
			},
		}
		Expect(k8sClient.Status().Update(ctx, deployment)).To(Succeed())
	}

	// Helper function to get a specific condition from Instance
	getCondition := func(instance *challengesv1.Instance, condType string) *metav1.Condition {
		for i, cond := range instance.Status.Conditions {
			if cond.Type == condType {
				return &instance.Status.Conditions[i]
			}
		}
		return nil
	}

	// Cleanup helper
	cleanupInstance := func(name, namespace string) {
		instance := &challengesv1.Instance{}
		err := k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, instance)
		if err == nil {
			_ = k8sClient.Delete(ctx, instance)
		}
	}

	It("should create Deployment and Service for new Instance", func() {
		instanceName := "instance-basic"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates resources")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Deployment was created")
		deployment := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{
				Name:      "deployment-" + instanceName,
				Namespace: namespace,
			}, deployment)
		}, timeout, interval).Should(Succeed())

		By("Verifying Service was created")
		service := &corev1.Service{}
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{
				Name:      "svc-" + instanceName,
				Namespace: namespace,
			}, service)
		}, timeout, interval).Should(Succeed())
	})

	It("should mark Instance Running when Deployment is ready", func() {
		instanceName := "instance-ready"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates resources")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Simulating Deployment becoming ready")
		simulateDeploymentReady(ctx, instanceName, namespace)

		By("Fourth reconcile - detects ready Deployment")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Phase is Running")
		updated := &challengesv1.Instance{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: instanceName, Namespace: namespace}, updated)).To(Succeed())
		Expect(updated.Status.Phase).To(Equal(challengesv1.PhaseRunning))
	})

	It("should remain Pending if Deployment fails", func() {
		instanceName := "instance-failed"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates resources")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Simulating Deployment failure")
		simulateDeploymentFailed(ctx, instanceName, namespace)

		By("Fourth reconcile - checks for pod failures")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Phase remains Pending (no pod failures detected)")
		updated := &challengesv1.Instance{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: instanceName, Namespace: namespace}, updated)).To(Succeed())
		Expect(updated.Status.Phase).To(Equal(challengesv1.PhasePending))
	})

	It("should set ServiceReady condition after Service creation", func() {
		instanceName := "instance-service"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates Service and sets condition")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying ServiceReady condition is set")
		updated := &challengesv1.Instance{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: instanceName, Namespace: namespace}, updated)).To(Succeed())

		cond := getCondition(updated, "ServiceReady")
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionTrue))
	})

	It("should delete Instance when expired", func() {
		instanceName := "instance-expired"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		past := metav1.NewTime(time.Now().Add(-time.Hour))
		instance := createBasicInstance(instanceName, namespace)
		instance.Spec.Lifecycle = &challengesv1.Lifecycle{ExpiresAt: &past}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - marks Instance as Expired")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Fourth reconcile - deletes expired Instance")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Instance is deleted or being deleted")
		Eventually(func() bool {
			updated := &challengesv1.Instance{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: instanceName, Namespace: namespace}, updated)
			return errors.IsNotFound(err) || !updated.DeletionTimestamp.IsZero()
		}, timeout, interval).Should(BeTrue())
	})

	It("should not create Service if no endpoints", func() {
		instanceName := "instance-no-svc"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		instance.Spec.Endpoints = nil
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - skips Service creation (no endpoints)")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Service was not created")
		service := &corev1.Service{}
		err = k8sClient.Get(ctx, types.NamespacedName{
			Name:      "svc-" + instanceName,
			Namespace: namespace,
		}, service)
		Expect(errors.IsNotFound(err)).To(BeTrue())
	})

	It("should set Instance to Staged phase when availableAt is in future", func() {
		instanceName := "instance-staged"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		future := metav1.NewTime(time.Now().Add(time.Hour))
		instance := createBasicInstance(instanceName, namespace)
		instance.Spec.Lifecycle = &challengesv1.Lifecycle{AvailableAt: &future}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates resources")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Simulating Deployment becoming ready")
		simulateDeploymentReady(ctx, instanceName, namespace)

		By("Fourth reconcile - should transition to Staged")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Phase is Staged (not Running)")
		updated := &challengesv1.Instance{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: instanceName, Namespace: namespace}, updated)).To(Succeed())
		Expect(updated.Status.Phase).To(Equal(challengesv1.PhaseStaged))
	})

	It("should update Deployment when spec changes", func() {
		instanceName := "instance-update-deploy"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates resources")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying initial Deployment")
		deployment := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{
				Name:      "deployment-" + instanceName,
				Namespace: namespace,
			}, deployment)
		}, timeout, interval).Should(Succeed())
		initialImage := deployment.Spec.Template.Spec.Containers[0].Image

		By("Updating Instance spec (this would normally be rejected by webhook)")
		// Note: In real scenario, the webhook prevents image changes
		// This test verifies controller behavior if spec somehow changes
		updated := &challengesv1.Instance{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: instanceName, Namespace: namespace}, updated)).To(Succeed())

		By("Adding resource limits to trigger reconciliation")
		updated.Spec.Limits = corev1.ResourceList{
			corev1.ResourceCPU:    *resourceQuantity("500m"),
			corev1.ResourceMemory: *resourceQuantity("512Mi"),
		}
		Expect(k8sClient.Update(ctx, updated)).To(Succeed())

		By("Fourth reconcile - detects spec change and updates Deployment")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Deployment was updated with new limits")
		Eventually(func() bool {
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      "deployment-" + instanceName,
				Namespace: namespace,
			}, deployment)
			if err != nil {
				return false
			}
			limits := deployment.Spec.Template.Spec.Containers[0].Resources.Limits
			return limits != nil && !limits.Cpu().IsZero()
		}, timeout, interval).Should(BeTrue())
		Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(Equal(initialImage))
	})

	It("should update Service when endpoints change", func() {
		instanceName := "instance-update-svc"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates Service with 1 port")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		service := &corev1.Service{}
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{
				Name:      "svc-" + instanceName,
				Namespace: namespace,
			}, service)
		}, timeout, interval).Should(Succeed())
		Expect(service.Spec.Ports).To(HaveLen(1))

		By("Adding another endpoint")
		updated := &challengesv1.Instance{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: instanceName, Namespace: namespace}, updated)).To(Succeed())
		updated.Spec.Endpoints = append(updated.Spec.Endpoints, challengesv1.EndpointSpec{
			Name:       "metrics",
			Protocol:   challengesv1.ProtocolHTTP,
			TargetPort: 9090,
		})
		Expect(k8sClient.Update(ctx, updated)).To(Succeed())

		By("Fourth reconcile - updates Service with 2 ports")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() int {
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      "svc-" + instanceName,
				Namespace: namespace,
			}, service)
			if err != nil {
				return 0
			}
			return len(service.Spec.Ports)
		}, timeout, interval).Should(Equal(2))
	})

	It("should set all conditions correctly", func() {
		instanceName := "instance-conditions"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates resources and sets conditions")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying conditions are set")
		updated := &challengesv1.Instance{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: instanceName, Namespace: namespace}, updated)).To(Succeed())

		deploymentCond := getCondition(updated, "DeploymentReady")
		Expect(deploymentCond).NotTo(BeNil())
		Expect(deploymentCond.Status).To(Equal(metav1.ConditionFalse))

		serviceCond := getCondition(updated, "ServiceReady")
		Expect(serviceCond).NotTo(BeNil())
		Expect(serviceCond.Status).To(Equal(metav1.ConditionTrue))

		ingressCond := getCondition(updated, "IngressReady")
		Expect(ingressCond).NotTo(BeNil())
		// IngressReady should be True since HTTP endpoint creates IngressRoute
		Expect(ingressCond.Status).To(Equal(metav1.ConditionTrue))
	})

	It("should handle on-demand challenge with team", func() {
		instanceName := "instance-ondemand"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		instance := createBasicInstance(instanceName, namespace)
		instance.Spec.Challenge.Type = challengesv1.ChallengeTypeOnDemand
		instance.Spec.Team = &challengesv1.Team{ID: 123}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - creates resources")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Deployment has team label")
		deployment := &appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{
				Name:      "deployment-" + instanceName,
				Namespace: namespace,
			}, deployment)
		}, timeout, interval).Should(Succeed())
		Expect(deployment.Labels).To(HaveKeyWithValue("challenges.isolet.dev/team", "123"))
	})

	It("should requeue Instance with future expiresAt", func() {
		instanceName := "instance-requeue"
		namespace := "default"
		defer cleanupInstance(instanceName, namespace)

		future := metav1.NewTime(time.Now().Add(5 * time.Second))
		instance := createBasicInstance(instanceName, namespace)
		instance.Spec.Lifecycle = &challengesv1.Lifecycle{ExpiresAt: &future}
		Expect(k8sClient.Create(ctx, instance)).To(Succeed())

		By("First reconcile - adds finalizer")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - sets Pending phase")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Third reconcile - reconciles resources and returns requeue duration")
		result, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.RequeueAfter).To(BeNumerically(">", 0))
		Expect(result.RequeueAfter).To(BeNumerically("<=", 5*time.Second))
	})
})

// Helper to create resource.Quantity for testing
func resourceQuantity(value string) *resource.Quantity {
	q := resource.MustParse(value)
	return &q
}

var _ = Describe("Instance Controller Unit Tests", func() {
	var reconciler *InstanceReconciler

	BeforeEach(func() {
		reconciler = &InstanceReconciler{
			Client:   k8sClient,
			Scheme:   k8sClient.Scheme(),
			Recorder: record.NewFakeRecorder(100),
		}
	})

	Describe("isDeploymentReady", func() {
		It("should return false when no replicas are ready", func() {
			deployment := &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Replicas:          1,
					ReadyReplicas:     0,
					AvailableReplicas: 0,
					UpdatedReplicas:   0,
				},
			}
			Expect(reconciler.isDeploymentReady(deployment)).To(BeFalse())
		})

		It("should return true when all replicas are ready (default 1)", func() {
			deployment := &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					Replicas:          1,
					ReadyReplicas:     1,
					AvailableReplicas: 1,
					UpdatedReplicas:   1,
				},
			}
			Expect(reconciler.isDeploymentReady(deployment)).To(BeTrue())
		})

		It("should use spec replicas when set", func() {
			replicas := int32(3)
			deployment := &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
				},
				Status: appsv1.DeploymentStatus{
					Replicas:          3,
					ReadyReplicas:     2,
					AvailableReplicas: 2,
					UpdatedReplicas:   2,
				},
			}
			Expect(reconciler.isDeploymentReady(deployment)).To(BeFalse())
		})

		It("should return true when all specified replicas are ready", func() {
			replicas := int32(2)
			deployment := &appsv1.Deployment{
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
				},
				Status: appsv1.DeploymentStatus{
					Replicas:          2,
					ReadyReplicas:     2,
					AvailableReplicas: 2,
					UpdatedReplicas:   2,
				},
			}
			Expect(reconciler.isDeploymentReady(deployment)).To(BeTrue())
		})
	})

	Describe("determinePhase", func() {
		It("should return Pending when deployment is not ready", func() {
			instance := &challengesv1.Instance{}
			phase := reconciler.determinePhase(instance, false, true, true)
			Expect(phase).To(Equal(challengesv1.PhasePending))
		})

		It("should return Running when all resources are ready and no lifecycle constraint", func() {
			instance := &challengesv1.Instance{}
			phase := reconciler.determinePhase(instance, true, true, true)
			Expect(phase).To(Equal(challengesv1.PhaseRunning))
		})

		It("should return Staged when deployment is ready but availableAt is in the future", func() {
			future := metav1.NewTime(time.Now().Add(time.Hour))
			instance := &challengesv1.Instance{
				Spec: challengesv1.InstanceSpec{
					Lifecycle: &challengesv1.Lifecycle{
						AvailableAt: &future,
					},
				},
			}
			phase := reconciler.determinePhase(instance, true, true, true)
			Expect(phase).To(Equal(challengesv1.PhaseStaged))
		})

		It("should return Running when availableAt is in the past", func() {
			past := metav1.NewTime(time.Now().Add(-time.Hour))
			instance := &challengesv1.Instance{
				Spec: challengesv1.InstanceSpec{
					Lifecycle: &challengesv1.Lifecycle{
						AvailableAt: &past,
					},
				},
			}
			phase := reconciler.determinePhase(instance, true, true, true)
			Expect(phase).To(Equal(challengesv1.PhaseRunning))
		})

		It("should return Pending when deployment is ready but service is not", func() {
			instance := &challengesv1.Instance{}
			phase := reconciler.determinePhase(instance, true, false, true)
			Expect(phase).To(Equal(challengesv1.PhasePending))
		})

		It("should return Pending when deployment is ready but ingress is not", func() {
			instance := &challengesv1.Instance{}
			phase := reconciler.determinePhase(instance, true, true, false)
			Expect(phase).To(Equal(challengesv1.PhasePending))
		})
	})

	Describe("isFailureReason", func() {
		DescribeTable("known failure reasons",
			func(reason string, expected bool) {
				Expect(isFailureReason(reason)).To(Equal(expected))
			},
			Entry("CrashLoopBackOff", "CrashLoopBackOff", true),
			Entry("ErrImagePull", "ErrImagePull", true),
			Entry("ImagePullBackOff", "ImagePullBackOff", true),
			Entry("CreateContainerConfigError", "CreateContainerConfigError", true),
			Entry("InvalidImageName", "InvalidImageName", true),
			Entry("CreateContainerError", "CreateContainerError", true),
			Entry("Running is not a failure", "Running", false),
			Entry("Completed is not a failure", "Completed", false),
			Entry("empty string is not a failure", "", false),
			Entry("OOMKilled is not a failure", "OOMKilled", false),
		)
	})

	Describe("equalIngressRouteSpec", func() {
		It("should return true for equal specs", func() {
			spec := &traefikv1alpha1.IngressRouteSpec{
				EntryPoints: []string{"web", "websecure"},
				Routes: []traefikv1alpha1.Route{
					{Kind: "Rule", Match: "Host(`example.com`)"},
				},
				TLS: &traefikv1alpha1.TLS{SecretName: "my-cert"},
			}
			Expect(equalIngressRouteSpec(spec, spec)).To(BeTrue())
		})

		It("should return false for different entry points count", func() {
			a := &traefikv1alpha1.IngressRouteSpec{
				EntryPoints: []string{"web"},
			}
			b := &traefikv1alpha1.IngressRouteSpec{
				EntryPoints: []string{"web", "websecure"},
			}
			Expect(equalIngressRouteSpec(a, b)).To(BeFalse())
		})

		It("should return false for different entry point values", func() {
			a := &traefikv1alpha1.IngressRouteSpec{
				EntryPoints: []string{"web"},
			}
			b := &traefikv1alpha1.IngressRouteSpec{
				EntryPoints: []string{"websecure"},
			}
			Expect(equalIngressRouteSpec(a, b)).To(BeFalse())
		})

		It("should return false for different route count", func() {
			a := &traefikv1alpha1.IngressRouteSpec{
				Routes: []traefikv1alpha1.Route{
					{Kind: "Rule", Match: "Host(`a.com`)"},
				},
			}
			b := &traefikv1alpha1.IngressRouteSpec{
				Routes: []traefikv1alpha1.Route{
					{Kind: "Rule", Match: "Host(`a.com`)"},
					{Kind: "Rule", Match: "Host(`b.com`)"},
				},
			}
			Expect(equalIngressRouteSpec(a, b)).To(BeFalse())
		})

		It("should return false when one spec has TLS and the other does not", func() {
			a := &traefikv1alpha1.IngressRouteSpec{
				TLS: &traefikv1alpha1.TLS{SecretName: "my-cert"},
			}
			b := &traefikv1alpha1.IngressRouteSpec{}
			Expect(equalIngressRouteSpec(a, b)).To(BeFalse())
		})

		It("should return false for different TLS secret names", func() {
			a := &traefikv1alpha1.IngressRouteSpec{
				TLS: &traefikv1alpha1.TLS{SecretName: "cert-a"},
			}
			b := &traefikv1alpha1.IngressRouteSpec{
				TLS: &traefikv1alpha1.TLS{SecretName: "cert-b"},
			}
			Expect(equalIngressRouteSpec(a, b)).To(BeFalse())
		})
	})

	Describe("equalEndpointStatus", func() {
		It("should return true for identical slices", func() {
			eps := []challengesv1.EndpointStatus{
				{
					EndpointSpec: challengesv1.EndpointSpec{
						Name:       "http",
						Protocol:   challengesv1.ProtocolHTTP,
						TargetPort: 80,
					},
					Hostname: "inst.challenge.isolet.dev",
					Ready:    true,
				},
			}
			Expect(equalEndpointStatus(eps, eps)).To(BeTrue())
		})

		It("should return true for empty slices", func() {
			Expect(equalEndpointStatus(nil, nil)).To(BeTrue())
			Expect(equalEndpointStatus([]challengesv1.EndpointStatus{}, []challengesv1.EndpointStatus{})).To(BeTrue())
		})

		It("should return false for different lengths", func() {
			a := []challengesv1.EndpointStatus{
				{EndpointSpec: challengesv1.EndpointSpec{Name: "http", Protocol: challengesv1.ProtocolHTTP, TargetPort: 80}},
			}
			Expect(equalEndpointStatus(a, nil)).To(BeFalse())
		})

		It("should return false for different hostnames", func() {
			a := []challengesv1.EndpointStatus{
				{
					EndpointSpec: challengesv1.EndpointSpec{Name: "http", Protocol: challengesv1.ProtocolHTTP, TargetPort: 80},
					Hostname:     "a.isolet.dev",
					Ready:        true,
				},
			}
			b := []challengesv1.EndpointStatus{
				{
					EndpointSpec: challengesv1.EndpointSpec{Name: "http", Protocol: challengesv1.ProtocolHTTP, TargetPort: 80},
					Hostname:     "b.isolet.dev",
					Ready:        true,
				},
			}
			Expect(equalEndpointStatus(a, b)).To(BeFalse())
		})

		It("should return false for different ready state", func() {
			a := []challengesv1.EndpointStatus{
				{
					EndpointSpec: challengesv1.EndpointSpec{Name: "http", Protocol: challengesv1.ProtocolHTTP, TargetPort: 80},
					Ready:        true,
				},
			}
			b := []challengesv1.EndpointStatus{
				{
					EndpointSpec: challengesv1.EndpointSpec{Name: "http", Protocol: challengesv1.ProtocolHTTP, TargetPort: 80},
					Ready:        false,
				},
			}
			Expect(equalEndpointStatus(a, b)).To(BeFalse())
		})

		It("should return false when endpoint name is missing in the second slice", func() {
			a := []challengesv1.EndpointStatus{
				{EndpointSpec: challengesv1.EndpointSpec{Name: "http", Protocol: challengesv1.ProtocolHTTP, TargetPort: 80}},
			}
			b := []challengesv1.EndpointStatus{
				{EndpointSpec: challengesv1.EndpointSpec{Name: "metrics", Protocol: challengesv1.ProtocolHTTP, TargetPort: 9090}},
			}
			Expect(equalEndpointStatus(a, b)).To(BeFalse())
		})
	})
})
