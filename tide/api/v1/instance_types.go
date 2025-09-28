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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Challenge type
// +kubebuilder:validation:Enum=dynamic;on-demand
type ChallengeType string

const (
	// ChallengeTypeDynamic indicates a challenge that is always running until manually terminated (common instance for all teams)
	ChallengeTypeDynamic ChallengeType = "dynamic"

	// ChallengeTypeOnDemand indicates a challenge that is provisioned on-demand and terminated after use (isolated instance for a team)
	ChallengeTypeOnDemand ChallengeType = "on-demand"
)

// Protocol indicates the protocol or service type of an endpoint
// +kubebuilder:validation:Enum=http;https;nc;ssh
type Protocol string

const (
	// ProtocolHTTP indicates HTTP protocol
	ProtocolHTTP Protocol = "http"

	// ProtocolHTTPS indicates HTTPS protocol
	ProtocolHTTPS Protocol = "https"

	// ProtocolNC indicates Netcat protocol
	ProtocolNC Protocol = "nc"

	// ProtocolSSH indicates SSH protocol
	ProtocolSSH Protocol = "ssh"
)

// Phase indicates the current phase of the instance
// Possible values:
// - Pending: Instance is being created
// - Staged: Instance is running but not accepting traffic yet (can be configured via lifecycle.availableAt)
// - Running: Instance is running and accepting traffic
// - Failed: Instance creation or startup failed
//
// +kubebuilder:validation:Enum=Pending;Staged;Running;Failed
type Phase string

const (
	// PhasePending indicates the instance is being created
	PhasePending Phase = "Pending"

	// PhaseStaged indicates the instance is running but not accepting traffic yet
	PhaseStaged Phase = "Staged"

	// PhaseRunning indicates the instance is running and accepting traffic
	PhaseRunning Phase = "Running"

	// PhaseFailed indicates the instance creation or startup failed
	PhaseFailed Phase = "Failed"
)

// Challenge object defines configuration for the challenge for which the instance is created
type Challenge struct {
	// ID of the challenge
	// +required
	// +kubebuilder:validation:Minimum=1
	ID int64 `json:"id"`

	// DNS-safe slug derived from challenge name
	// ex: Alice's Adventure -> alices-adventure
	// +required
	Name string `json:"name"`

	// A custom and unique flag can be passed per instance
	// If set, it will be injected as an environment variable into the container
	// +optional
	Flag *string `json:"flag,omitempty"`

	// type indicates the provisioning type of the challenge
	// +required
	Type ChallengeType `json:"type"`

	// Docker image to use for the challenge instance
	// +required
	Image string `json:"image"`
}

// Team refers to the team for which the instance is created
type Team struct {
	// ID of the team
	// +required
	// +kubebuilder:validation:Minimum=1
	ID int64 `json:"id"`
}

// Endpoint represents a network endpoint for an instance
type Endpoint struct {
	// Name of the endpoint (e.g., "app", "metrics")
	// +optional
	Name *string `json:"name,omitempty"`

	// Protocol of the endpoint (e.g., "http", "ssh")
	// +required
	Protocol Protocol `json:"protocol"`

	// Hostname of the endpoint
	// Controller genarates a DNS-safe hostname if not provided
	// e.g., <instance-uuid>.<challenge-name>.<tide-domain>
	//
	// Example: "a1b2c3d4-e5f6-7890-abcd-ef1234567890.alices-adventure.isolet.dev"
	// +optional
	Hostname *string `json:"hostname,omitempty"`

	// Port of the endpoint
	// This is resolved by the controller based on the service configuration
	// +optional
	Port *int32 `json:"port"`

	// TargetPort of the endpoint
	// Port on which the application inside the container is listening
	// +required
	TargetPort int32 `json:"targetPort"`

	// Ready indicates whether the endpoint is ready to accept traffic
	// +optional
	Ready *bool `json:"ready"`
}

// Lifecycle defines the lifecycle configuration for an instance
type Lifecycle struct {
	// timestamp of the time when the instance should be available
	// If not set, the instance will be available immediately after provisioning
	// +optional
	AvailableAt *metav1.Time `json:"availableAt,omitempty"`

	// timestamp of the time when the instance should be terminated
	// If not set, the instance will run indefinitely until manually terminated
	// +optional
	ExpiresAt *metav1.Time `json:"expiresAt,omitempty"`

	// allowExtension indicates whether the instance's expiry time can be extended by the user
	// Default is true
	// +optional
	// +kubebuilder:default=true
	AllowExtension bool `json:"allowExtension"`

	// restartPolicy defines the restart policy for all containers within the instance
	// One of Always, OnFailure, Never
	// Default is Never
	// +optional
	// +kubebuilder:default="Never"
	RestartPolicy corev1.RestartPolicy `json:"restartPolicy,omitempty"`
}

// InstanceSpec defines the desired state of Instance
type InstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// The following markers will use OpenAPI v3 schema to validate the value
	// More info: https://book.kubebuilder.io/reference/markers/crd-validation.html

	// challenge is the challenge for which this instance is created
	// +required
	Challenge Challenge `json:"challenge"`

	// team is the team for which this instance is created
	// dynamic challenges may not have a team associated
	// +optional
	Team *Team `json:"team,omitempty"`

	// requests is the resource requests for the instance
	// +optional
	Requests corev1.ResourceList `json:"requests,omitempty"`

	// limits is the resource limits for the instance
	// +optional
	Limits corev1.ResourceList `json:"limits,omitempty"`

	// Endpoints represents the network endpoints for the instance
	// +optional
	Endpoints []Endpoint `json:"endpoints,omitempty"`

	// lifecycle defines the lifecycle configuration for the instance
	// +optional
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
}

// InstanceStatus defines the observed state of Instance.
type InstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// phase indicates the current phase of the instance
	// +optional
	Phase Phase `json:"phase,omitempty"` // Pending | Staged | Running | Failed

	// endpoints are the resolved network endpoints for the instance
	// +optional
	Endpoints []Endpoint `json:"endpoints,omitempty"`

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the Instance resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Instance is the Schema for the instances API
type Instance struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of Instance
	// +required
	Spec InstanceSpec `json:"spec"`

	// status defines the observed state of Instance
	// +optional
	Status InstanceStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// InstanceList contains a list of Instance
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Instance{}, &InstanceList{})
}
