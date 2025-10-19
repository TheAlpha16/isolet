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

// ChallengeType represents how a challenge is provisioned.
// +kubebuilder:validation:Enum=dynamic;on-demand
type ChallengeType string

const (
	// ChallengeTypeDynamic: a long-running, shared instance available to all teams.
	ChallengeTypeDynamic ChallengeType = "dynamic"

	// ChallengeTypeOnDemand: a short-lived, isolated instance created for a specific team.
	ChallengeTypeOnDemand ChallengeType = "on-demand"
)

// Protocol defines the network protocol or service type of an endpoint.
// +kubebuilder:validation:Enum=http;https;nc;ssh
type Protocol string

const (
	ProtocolHTTP  Protocol = "http"  // HTTP service
	ProtocolHTTPS Protocol = "https" // HTTPS service
	ProtocolNC    Protocol = "nc"    // Netcat-style TCP service
	ProtocolSSH   Protocol = "ssh"   // SSH service
)

// Phase describes the current lifecycle phase of an Instance.
// +kubebuilder:validation:Enum=Pending;Staged;Running;Failed
type Phase string

const (
	// Instance has been created but resources are not yet ready.
	PhasePending Phase = "Pending"

	// Instance is running but not yet available (e.g., waiting for availableAt).
	PhaseStaged Phase = "Staged"

	// Instance is fully available and accepting traffic.
	PhaseRunning Phase = "Running"

	// Instance failed during creation or startup.
	PhaseFailed Phase = "Failed"
)

// Challenge contains metadata about the challenge backing an Instance.
type Challenge struct {
	// Unique ID of the challenge.
	// +required
	// +kubebuilder:validation:Minimum=1
	ID int64 `json:"id"`

	// DNS-safe slug derived from the challenge name.
	// Example: "Alice's Adventure" -> "alices-adventure".
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	// +required
	Slug string `json:"slug"`

	// Optional flag injected into the instance as an environment variable.
	// +optional
	Flag *string `json:"flag,omitempty"`

	// How this challenge is provisioned (dynamic or on-demand).
	// +required
	Type ChallengeType `json:"type"`

	// Container image used to run the challenge.
	// +required
	Image string `json:"image"`
}

// Team identifies the team an Instance belongs to.
type Team struct {
	// Unique ID of the team.
	// +required
	// +kubebuilder:validation:Minimum=1
	ID int64 `json:"id"`
}

// EndpointSpec defines a single network entrypoint for an Instance.
type EndpointSpec struct {
	// Logical name for the endpoint (e.g., "app", "metrics").
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name"`

	// Network protocol for this endpoint (e.g., http, ssh).
	// +required
	Protocol Protocol `json:"protocol"`

	// Port on which the application inside the container is listening.
	// +required
	TargetPort int32 `json:"targetPort"`
}

// EndpointStatus defines a resolved network entrypoint for an Instance.
type EndpointStatus struct {
	// Embeds the spec fields.
	EndpointSpec `json:",inline"`

	// Fully-qualified domain name. If not provided, the controller will generate one.
	// Example: "<instance-uuid>.<challenge-name>.isolet.dev".
	// +optional
	Hostname *string `json:"hostname,omitempty"`

	// Exposed service port assigned by controller (cluster-facing).
	// +optional
	Port *int32 `json:"port"`

	// Whether this endpoint is ready to receive traffic.
	// +optional
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
}

// Lifecycle configures availability and expiry of an Instance.
type Lifecycle struct {
	// Time when the Instance should first become available.
	// If unset, the instance is available immediately after provisioning.
	// +optional
	AvailableAt *metav1.Time `json:"availableAt,omitempty"`

	// Time when the Instance should be terminated.
	// If unset, the instance will run until manually deleted.
	// +optional
	ExpiresAt *metav1.Time `json:"expiresAt,omitempty"`

	// Whether expiry can be extended by the user (default: true).
	// +optional
	// +kubebuilder:default=true
	AllowExtension bool `json:"allowExtension"`
}

func (l *Lifecycle) SetDefaults() {
	l.AllowExtension = true
}

// InstanceSpec defines the desired state of the Instance.
type InstanceSpec struct {
	// Challenge backing this instance.
	// +required
	Challenge Challenge `json:"challenge"`

	// Team this instance belongs to (not required for dynamic challenges).
	// +optional
	Team *Team `json:"team,omitempty"`

	// Resource requests for CPU/memory.
	// +optional
	Requests corev1.ResourceList `json:"requests,omitempty"`

	// Resource limits for CPU/memory.
	// +optional
	Limits corev1.ResourceList `json:"limits,omitempty"`

	// Network endpoints exposed by this instance.
	// +listType=map
	// +listMapKey=name
	// +optional
	Endpoints []EndpointSpec `json:"endpoints,omitempty"`

	// Lifecycle configuration (availability and expiry).
	// +optional
	Lifecycle *Lifecycle `json:"lifecycle,omitempty"`
}

// InstanceStatus captures the observed state of an Instance.
type InstanceStatus struct {
	// Current lifecycle phase.
	// +optional
	Phase Phase `json:"phase,omitempty"`

	// Resolved endpoints after provisioning.
	// +listType=map
	// +listMapKey=name
	// +optional
	Endpoints []EndpointStatus `json:"endpoints,omitempty"`

	// Standard set of conditions describing resource health and transitions.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Instance is the Schema for the instances API.
// +kubebuilder:resource:shortName=inst
// +kubebuilder:printcolumn:name="Challenge",type=string,JSONPath=`.spec.challenge.name`,description="Challenge Name"
// +kubebuilder:printcolumn:name="Team",type=string,JSONPath=`.spec.team.id`,description="Team ID"
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.spec.challenge.type`,description="Challenge Type"
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`,description="Instance Phase",priority=1
// +kubebuilder:printcolumn:name="ExpiresAt",type=string,JSONPath=`.spec.lifecycle.expiresAt`,description="Instance expiry time"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="Time since creation"
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// Desired configuration.
	Spec InstanceSpec `json:"spec"`

	// Current status.
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
