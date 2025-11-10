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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
					ID:    1,
					Slug:  "test-challenge",
					Type:  challengesv1.ChallengeTypeDynamic,
					Image: "nginx:latest",
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

		By("First reconcile - sets Pending phase")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - creates resources")
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

		By("First reconcile - sets Pending")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - creates resources")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Simulating Deployment becoming ready")
		simulateDeploymentReady(ctx, instanceName, namespace)

		By("Third reconcile - detects ready Deployment")
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

		By("First reconcile - sets Pending")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - creates resources")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Simulating Deployment failure")
		simulateDeploymentFailed(ctx, instanceName, namespace)

		By("Third reconcile - detects failed Deployment")
		_, err = reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Verifying Phase remains Pending")
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

		By("First reconcile - sets Pending")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - creates Service and sets condition")
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

		By("First reconcile - sets Pending phase")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - checks expiry and deletes")
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

		By("First reconcile - sets Pending")
		_, err := reconciler.Reconcile(ctx, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: instanceName, Namespace: namespace},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Second reconcile - skips Service creation")
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
})
