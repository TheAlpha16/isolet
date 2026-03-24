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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	challengesv1 "github.com/TheAlpha16/isolet/tide/api/v1"
)

var _ = Describe("Instance Webhook", func() {
	var (
		validator InstanceCustomValidator
		defaulter InstanceCustomDefaulter
	)

	// Helper function to create a valid instance
	createValidInstance := func() *challengesv1.Instance {
		return &challengesv1.Instance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-instance",
				Namespace: "default",
			},
			Spec: challengesv1.InstanceSpec{
				Challenge: challengesv1.Challenge{
					ID:     1,
					Slug:   "test-challenge",
					Type:   challengesv1.ChallengeTypeDynamic,
					Image:  "nginx:latest",
					Domain: "isolet.dev",
				},
			},
		}
	}

	BeforeEach(func() {
		validator = InstanceCustomValidator{}
		defaulter = InstanceCustomDefaulter{}
		Expect(validator).NotTo(BeNil(), "Expected validator to be initialized")
		Expect(defaulter).NotTo(BeNil(), "Expected defaulter to be initialized")
	})

	Context("When creating Instance under Defaulting Webhook", func() {
		It("Should apply default lifecycle when lifecycle is nil", func() {
			By("Creating an instance without lifecycle")
			obj := createValidInstance()
			obj.Spec.Lifecycle = nil

			By("Calling the Default method to apply defaults")
			err := defaulter.Default(ctx, obj)

			By("Checking that defaults are applied")
			Expect(err).NotTo(HaveOccurred())
			Expect(obj.Spec.Lifecycle).NotTo(BeNil())
			Expect(obj.Spec.Lifecycle.AllowExtension).To(BeTrue())
		})

		It("Should not override existing lifecycle values", func() {
			By("Creating an instance with custom lifecycle values")
			futureTime := metav1.NewTime(time.Now().Add(1 * time.Hour))
			obj := createValidInstance()
			obj.Spec.Lifecycle = &challengesv1.Lifecycle{
				AllowExtension: false,
				ExpiresAt:      &futureTime,
			}

			By("Calling the Default method")
			err := defaulter.Default(ctx, obj)

			By("Checking that custom values are preserved")
			Expect(err).NotTo(HaveOccurred())
			Expect(obj.Spec.Lifecycle.AllowExtension).To(BeFalse())
			Expect(obj.Spec.Lifecycle.ExpiresAt).To(Equal(&futureTime))
		})

		It("Should return error for wrong object type", func() {
			By("Passing a wrong object type to defaulter")
			wrongObj := &corev1.Pod{}

			By("Calling the Default method")
			err := defaulter.Default(ctx, wrongObj)

			By("Checking that an error is returned")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected an Instance object"))
		})
	})

	Context("When creating Instance under Validating Webhook", func() {
		It("Should reject missing domain", func() {
			By("Creating an instance without a domain")
			obj := createValidInstance()
			obj.Spec.Challenge.Domain = ""

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails for missing domain")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("challenge.domain is required"))
			Expect(warnings).To(BeNil())
		})

		It("Should accept valid dynamic challenge instance", func() {
			By("Creating a valid dynamic challenge instance")
			obj := createValidInstance()

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation passes")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should accept valid on-demand challenge with team", func() {
			By("Creating a valid on-demand challenge with team specified")
			obj := createValidInstance()
			obj.Spec.Challenge.Type = challengesv1.ChallengeTypeOnDemand
			obj.Spec.Team = &challengesv1.Team{ID: 1}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation passes")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should reject on-demand challenge without team", func() {
			By("Creating an on-demand challenge without team")
			obj := createValidInstance()
			obj.Spec.Challenge.Type = challengesv1.ChallengeTypeOnDemand
			obj.Spec.Team = nil

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails with appropriate error")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("team must be set for on-demand challenges"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject duplicate endpoint names", func() {
			By("Creating an instance with duplicate endpoint names")
			obj := createValidInstance()
			obj.Spec.Endpoints = []challengesv1.EndpointSpec{
				{Name: "web", Protocol: challengesv1.ProtocolHTTP, TargetPort: 80},
				{Name: "web", Protocol: challengesv1.ProtocolHTTPS, TargetPort: 443},
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails for duplicate names")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate endpoint name: web"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject duplicate endpoint target ports", func() {
			By("Creating an instance with duplicate target ports")
			obj := createValidInstance()
			obj.Spec.Endpoints = []challengesv1.EndpointSpec{
				{Name: "http", Protocol: challengesv1.ProtocolHTTP, TargetPort: 8080},
				{Name: "alt", Protocol: challengesv1.ProtocolHTTPS, TargetPort: 8080},
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails for duplicate target ports")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate endpoint targetPort: 8080"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject invalid endpoint port (below range)", func() {
			By("Creating an instance with port below valid range")
			obj := createValidInstance()
			obj.Spec.Endpoints = []challengesv1.EndpointSpec{
				{Name: "web", Protocol: challengesv1.ProtocolHTTP, TargetPort: 0},
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails for invalid port")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("endpoint port must be between 1 and 65535"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject port above valid range", func() {
			By("Creating an instance with port above 65535")
			obj := createValidInstance()
			obj.Spec.Endpoints = []challengesv1.EndpointSpec{
				{Name: "web", Protocol: challengesv1.ProtocolHTTP, TargetPort: 65536},
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails for port out of range")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("endpoint port must be between 1 and 65535"))
			Expect(warnings).To(BeNil())
		})

		It("Should accept all valid protocols", func() {
			protocols := []challengesv1.Protocol{
				challengesv1.ProtocolHTTP,
				challengesv1.ProtocolHTTPS,
				challengesv1.ProtocolNC,
				challengesv1.ProtocolSSH,
			}

			for _, proto := range protocols {
				By("Creating instance with protocol: " + string(proto))
				obj := createValidInstance()
				obj.Spec.Endpoints = []challengesv1.EndpointSpec{
					{Name: "endpoint", Protocol: proto, TargetPort: 8080},
				}

				By("Validating the instance")
				warnings, err := validator.ValidateCreate(ctx, obj)

				By("Checking that validation passes for valid protocol")
				Expect(err).NotTo(HaveOccurred())
				Expect(warnings).To(BeNil())
			}
		})

		It("Should reject invalid protocol", func() {
			By("Creating an instance with invalid protocol")
			obj := createValidInstance()
			obj.Spec.Endpoints = []challengesv1.EndpointSpec{
				{Name: "web", Protocol: challengesv1.Protocol("ftp"), TargetPort: 21},
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails for invalid protocol")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("endpoint protocol must be one of"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject expiry in the past", func() {
			By("Creating an instance with past expiry time")
			pastTime := metav1.NewTime(time.Now().Add(-1 * time.Hour))
			obj := createValidInstance()
			obj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt: &pastTime,
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails for past expiry")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("lifecycle.expiresAt must be in the future"))
			Expect(warnings).To(BeNil())
		})

		It("Should accept future expiry", func() {
			By("Creating an instance with future expiry time")
			futureTime := metav1.NewTime(time.Now().Add(1 * time.Hour))
			obj := createValidInstance()
			obj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt: &futureTime,
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation passes for future expiry")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should reject availableAt after expiresAt", func() {
			By("Creating an instance with availableAt after expiresAt")
			now := time.Now()
			availableTime := metav1.NewTime(now.Add(2 * time.Hour))
			expiresTime := metav1.NewTime(now.Add(1 * time.Hour))
			obj := createValidInstance()
			obj.Spec.Lifecycle = &challengesv1.Lifecycle{
				AvailableAt: &availableTime,
				ExpiresAt:   &expiresTime,
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation fails for invalid time ordering")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("lifecycle.availableAt must be before lifecycle.expiresAt"))
			Expect(warnings).To(BeNil())
		})

		It("Should accept valid availableAt and expiresAt", func() {
			By("Creating an instance with valid time ordering")
			now := time.Now()
			availableTime := metav1.NewTime(now.Add(1 * time.Hour))
			expiresTime := metav1.NewTime(now.Add(2 * time.Hour))
			obj := createValidInstance()
			obj.Spec.Lifecycle = &challengesv1.Lifecycle{
				AvailableAt: &availableTime,
				ExpiresAt:   &expiresTime,
			}

			By("Validating the instance")
			warnings, err := validator.ValidateCreate(ctx, obj)

			By("Checking that validation passes")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should return error for wrong object type", func() {
			By("Passing a wrong object type to validator")
			wrongObj := &corev1.Pod{}

			By("Validating the wrong object")
			warnings, err := validator.ValidateCreate(ctx, wrongObj)

			By("Checking that an error is returned")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected a Instance object"))
			Expect(warnings).To(BeNil())
		})
	})

	Context("When updating Instance under Validating Webhook", func() {
		It("Should accept valid update with no immutable changes", func() {
			By("Creating old and new instances with same immutable fields")
			oldObj := createValidInstance()
			oldObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				AllowExtension: true,
			}
			newObj := createValidInstance()
			futureTime := metav1.NewTime(time.Now().Add(1 * time.Hour))
			newObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &futureTime,
				AllowExtension: true,
			}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation passes")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should reject challenge.id change", func() {
			By("Creating instances with different challenge IDs")
			oldObj := createValidInstance()
			newObj := createValidInstance()
			newObj.Spec.Challenge.ID = 2

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for immutable field")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("challenge.id is immutable"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject challenge.slug change", func() {
			By("Creating instances with different challenge slugs")
			oldObj := createValidInstance()
			newObj := createValidInstance()
			newObj.Spec.Challenge.Slug = "different-challenge"

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for immutable field")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("challenge.slug is immutable"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject challenge.image change", func() {
			By("Creating instances with different challenge images")
			oldObj := createValidInstance()
			newObj := createValidInstance()
			newObj.Spec.Challenge.Image = "alpine:latest"

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for immutable field")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("challenge.image is immutable"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject challenge.domain change", func() {
			By("Creating instances with different challenge domains")
			oldObj := createValidInstance()
			newObj := createValidInstance()
			newObj.Spec.Challenge.Domain = "other.dev"

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for immutable domain")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("challenge.domain is immutable"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject challenge.type change", func() {
			By("Creating instances with different challenge types")
			oldObj := createValidInstance()
			oldObj.Spec.Challenge.Type = challengesv1.ChallengeTypeDynamic
			newObj := createValidInstance()
			newObj.Spec.Challenge.Type = challengesv1.ChallengeTypeOnDemand
			newObj.Spec.Team = &challengesv1.Team{ID: 1}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for immutable field")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("challenge.type is immutable"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject team change when team was set", func() {
			By("Creating instances with different teams")
			oldObj := createValidInstance()
			oldObj.Spec.Team = &challengesv1.Team{ID: 1}
			newObj := createValidInstance()
			newObj.Spec.Team = &challengesv1.Team{ID: 2}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for immutable team")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("team is immutable"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject team removal", func() {
			By("Creating instances where team is removed")
			oldObj := createValidInstance()
			oldObj.Spec.Team = &challengesv1.Team{ID: 1}
			newObj := createValidInstance()
			newObj.Spec.Team = nil

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for team removal")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("team is immutable"))
			Expect(warnings).To(BeNil())
		})

		It("Should reject team addition", func() {
			By("Creating instances where team is added")
			oldObj := createValidInstance()
			oldObj.Spec.Team = nil
			newObj := createValidInstance()
			newObj.Spec.Team = &challengesv1.Team{ID: 1}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for team addition")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("team is immutable"))
			Expect(warnings).To(BeNil())
		})

		It("Should accept expiry extension when allowExtension is true", func() {
			By("Creating instances with extended expiry time")
			now := time.Now()
			oldExpiry := metav1.NewTime(now.Add(1 * time.Hour))
			newExpiry := metav1.NewTime(now.Add(2 * time.Hour))

			oldObj := createValidInstance()
			oldObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &oldExpiry,
				AllowExtension: true,
			}

			newObj := createValidInstance()
			newObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &newExpiry,
				AllowExtension: true,
			}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation passes when extension is allowed")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should reject expiry extension when allowExtension is false", func() {
			By("Creating instances with extended expiry but extension disabled")
			now := time.Now()
			oldExpiry := metav1.NewTime(now.Add(1 * time.Hour))
			newExpiry := metav1.NewTime(now.Add(2 * time.Hour))

			oldObj := createValidInstance()
			oldObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &oldExpiry,
				AllowExtension: false,
			}

			newObj := createValidInstance()
			newObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &newExpiry,
				AllowExtension: false,
			}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails when extension is not allowed")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("lifecycle.expiresAt cannot be changed"))
			Expect(warnings).To(BeNil())
		})

		It("Should accept removing expiry", func() {
			By("Creating instances where expiry is removed")
			oldExpiry := metav1.NewTime(time.Now().Add(1 * time.Hour))

			oldObj := createValidInstance()
			oldObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &oldExpiry,
				AllowExtension: false,
			}

			newObj := createValidInstance()
			newObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      nil,
				AllowExtension: false,
			}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation passes for expiry removal")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should accept adding expiry when none existed", func() {
			By("Creating instances where expiry is added")
			newExpiry := metav1.NewTime(time.Now().Add(1 * time.Hour))

			oldObj := createValidInstance()
			oldObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				AllowExtension: true,
			}

			newObj := createValidInstance()
			newObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &newExpiry,
				AllowExtension: true,
			}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation passes for adding expiry")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should accept no change to expiry", func() {
			By("Creating instances with same expiry time")
			expiry := metav1.NewTime(time.Now().Add(1 * time.Hour))

			oldObj := createValidInstance()
			oldObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &expiry,
				AllowExtension: false,
			}

			newObj := createValidInstance()
			newObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt:      &expiry,
				AllowExtension: false,
			}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation passes when no change occurs")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})

		It("Should validate new endpoint configuration", func() {
			By("Creating update with invalid endpoint")
			oldObj := createValidInstance()
			newObj := createValidInstance()
			newObj.Spec.Endpoints = []challengesv1.EndpointSpec{
				{Name: "web", Protocol: challengesv1.ProtocolHTTP, TargetPort: 0},
			}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for invalid endpoint")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("endpoint port must be between 1 and 65535"))
			Expect(warnings).To(BeNil())
		})

		It("Should validate new lifecycle configuration", func() {
			By("Creating update with invalid lifecycle")
			oldObj := createValidInstance()
			newObj := createValidInstance()
			pastTime := metav1.NewTime(time.Now().Add(-1 * time.Hour))
			newObj.Spec.Lifecycle = &challengesv1.Lifecycle{
				ExpiresAt: &pastTime,
			}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, newObj)

			By("Checking that validation fails for past expiry")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("lifecycle.expiresAt must be in the future"))
			Expect(warnings).To(BeNil())
		})

		It("Should return error for wrong oldObj type", func() {
			By("Passing a wrong oldObj type")
			wrongObj := &corev1.Pod{}
			newObj := createValidInstance()

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, wrongObj, newObj)

			By("Checking that an error is returned")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected a Instance object for the oldObj"))
			Expect(warnings).To(BeNil())
		})

		It("Should return error for wrong newObj type", func() {
			By("Passing a wrong newObj type")
			oldObj := createValidInstance()
			wrongObj := &corev1.Pod{}

			By("Validating the update")
			warnings, err := validator.ValidateUpdate(ctx, oldObj, wrongObj)

			By("Checking that an error is returned")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("expected a Instance object for the newObj"))
			Expect(warnings).To(BeNil())
		})
	})

	Context("When deleting Instance under Validating Webhook", func() {
		It("Should always allow deletion", func() {
			By("Creating an instance to delete")
			obj := createValidInstance()

			By("Validating the deletion")
			warnings, err := validator.ValidateDelete(ctx, obj)

			By("Checking that validation passes")
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeNil())
		})
	})
})
