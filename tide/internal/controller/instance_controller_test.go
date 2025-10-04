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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	challengesv1 "github.com/TheAlpha16/isolet/tide/api/v1"
)

var _ = Describe("Instance Controller", func() {
	const (
		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When reconciling a new Instance", func() {
		const resourceName = "test-instance-new"
		const resourceNamespace = "default"

		ctx := context.Background()
		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: resourceNamespace,
		}

		AfterEach(func() {
			resource := &challengesv1.Instance{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				By("Cleanup the Instance resource")
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("Should eventually reach Running phase after reconciliation", func() {
			By("Creating a new Instance without phase")
			instance := &challengesv1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: resourceNamespace,
				},
				Spec: challengesv1.InstanceSpec{
					Challenge: challengesv1.Challenge{
						ID:    1,
						Name:  "test-challenge",
						Type:  challengesv1.ChallengeTypeDynamic,
						Image: "nginx:latest",
					},
				},
			}
			Expect(k8sClient.Create(ctx, instance)).To(Succeed())

			By("Reconciling the Instance twice to reach Running phase")
			controllerReconciler := &InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			// First reconcile - sets to Pending
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Second reconcile - transitions to Running
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeZero())

			By("Checking that Phase is Running")
			updatedInstance := &challengesv1.Instance{}
			err = k8sClient.Get(ctx, typeNamespacedName, updatedInstance)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedInstance.Status.Phase).To(Equal(challengesv1.PhaseRunning))
		})
	})

	Context("When reconciling Instance with expiry", func() {
		const resourceName = "test-instance-expiry"
		const resourceNamespace = "default"

		ctx := context.Background()
		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: resourceNamespace,
		}

		AfterEach(func() {
			resource := &challengesv1.Instance{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil && resource.DeletionTimestamp.IsZero() {
				By("Cleanup the Instance resource")
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("Should delete Instance when ExpiresAt is in the past", func() {
			By("Creating an Instance with past expiry")
			pastTime := metav1.NewTime(time.Now().Add(-1 * time.Hour))
			instance := &challengesv1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: resourceNamespace,
				},
				Spec: challengesv1.InstanceSpec{
					Challenge: challengesv1.Challenge{
						ID:    1,
						Name:  "test-challenge",
						Type:  challengesv1.ChallengeTypeDynamic,
						Image: "nginx:latest",
					},
					Lifecycle: &challengesv1.Lifecycle{
						ExpiresAt: &pastTime,
					},
				},
			}
			Expect(k8sClient.Create(ctx, instance)).To(Succeed())

			By("Setting Phase to Running")
			instance.Status.Phase = challengesv1.PhaseRunning
			Expect(k8sClient.Status().Update(ctx, instance)).To(Succeed())

			By("Reconciling the Instance")
			controllerReconciler := &InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeZero())

			By("Checking that Instance was deleted")
			deletedInstance := &challengesv1.Instance{}
			err = k8sClient.Get(ctx, typeNamespacedName, deletedInstance)
			Expect(errors.IsNotFound(err) || !deletedInstance.DeletionTimestamp.IsZero()).To(BeTrue())
		})

		It("Should requeue Instance when ExpiresAt is in the future", func() {
			By("Creating an Instance with future expiry")
			futureTime := metav1.NewTime(time.Now().Add(1 * time.Hour))
			instance := &challengesv1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: resourceNamespace,
				},
				Spec: challengesv1.InstanceSpec{
					Challenge: challengesv1.Challenge{
						ID:    1,
						Name:  "test-challenge",
						Type:  challengesv1.ChallengeTypeDynamic,
						Image: "nginx:latest",
					},
					Lifecycle: &challengesv1.Lifecycle{
						ExpiresAt: &futureTime,
					},
				},
			}
			Expect(k8sClient.Create(ctx, instance)).To(Succeed())

			By("Setting Phase to Running")
			instance.Status.Phase = challengesv1.PhaseRunning
			Expect(k8sClient.Status().Update(ctx, instance)).To(Succeed())

			By("Reconciling the Instance")
			controllerReconciler := &InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Checking that reconcile returns requeue with appropriate delay")
			Expect(result.RequeueAfter).To(BeNumerically(">", 0))
			Expect(result.RequeueAfter).To(BeNumerically("<=", 1*time.Hour))
		})

		It("Should not requeue Instance without expiry", func() {
			By("Creating an Instance without lifecycle expiry")
			instance := &challengesv1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: resourceNamespace,
				},
				Spec: challengesv1.InstanceSpec{
					Challenge: challengesv1.Challenge{
						ID:    1,
						Name:  "test-challenge",
						Type:  challengesv1.ChallengeTypeDynamic,
						Image: "nginx:latest",
					},
				},
				Status: challengesv1.InstanceStatus{
					Phase: challengesv1.PhaseRunning,
				},
			}
			Expect(k8sClient.Create(ctx, instance)).To(Succeed())

			By("Reconciling the Instance")
			controllerReconciler := &InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Checking that reconcile does not requeue")
			Expect(result.RequeueAfter).To(BeZero())
		})
	})

	Context("When reconciling deleted Instance", func() {
		const resourceName = "test-instance-delete"
		const resourceNamespace = "default"

		ctx := context.Background()
		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: resourceNamespace,
		}

		It("Should handle non-existent Instance gracefully", func() {
			By("Reconciling a non-existent Instance")
			controllerReconciler := &InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})

			By("Checking that reconcile completes without error")
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeZero())
		})

		It("Should handle Instance with DeletionTimestamp", func() {
			By("Creating an Instance")
			instance := &challengesv1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: resourceNamespace,
				},
				Spec: challengesv1.InstanceSpec{
					Challenge: challengesv1.Challenge{
						ID:    1,
						Name:  "test-challenge",
						Type:  challengesv1.ChallengeTypeDynamic,
						Image: "nginx:latest",
					},
				},
				Status: challengesv1.InstanceStatus{
					Phase: challengesv1.PhaseRunning,
				},
			}
			Expect(k8sClient.Create(ctx, instance)).To(Succeed())

			By("Deleting the Instance")
			Expect(k8sClient.Delete(ctx, instance)).To(Succeed())

			By("Reconciling the Instance with DeletionTimestamp")
			controllerReconciler := &InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})

			By("Checking that reconcile completes without error")
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeZero())
		})
	})

	Context("When reconciling Instance in Running phase", func() {
		const resourceName = "test-instance-running"
		const resourceNamespace = "default"

		ctx := context.Background()
		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: resourceNamespace,
		}

		AfterEach(func() {
			resource := &challengesv1.Instance{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				By("Cleanup the Instance resource")
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}
		})

		It("Should remain in Running phase if already Running", func() {
			By("Creating an Instance in Running phase")
			instance := &challengesv1.Instance{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: resourceNamespace,
				},
				Spec: challengesv1.InstanceSpec{
					Challenge: challengesv1.Challenge{
						ID:    1,
						Name:  "test-challenge",
						Type:  challengesv1.ChallengeTypeDynamic,
						Image: "nginx:latest",
					},
				},
			}
			Expect(k8sClient.Create(ctx, instance)).To(Succeed())

			By("Setting Phase to Running")
			instance.Status.Phase = challengesv1.PhaseRunning
			Expect(k8sClient.Status().Update(ctx, instance)).To(Succeed())

			By("Reconciling the Instance")
			controllerReconciler := &InstanceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			result, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.RequeueAfter).To(BeZero())

			By("Checking that Phase remains Running")
			updatedInstance := &challengesv1.Instance{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, updatedInstance)).To(Succeed())
			Expect(updatedInstance.Status.Phase).To(Equal(challengesv1.PhaseRunning))
		})
	})
})
