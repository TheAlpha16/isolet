package instance

import (
	"context"
	"fmt"
	"time"

	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	"github.com/TheAlpha16/isolet/api/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type instanceService struct {
	dynamicClient dynamic.Interface
	namespace     string
	config        *utils.Config
	
	// Fields set before each operation
	teamID      int64
	challengeID int64
	challenge   *ChallengeInfo
	instance    *instanceDom.Instance
}

// ChallengeInfo contains challenge metadata needed for instance creation
type ChallengeInfo struct {
	ID    int64
	Name  string
	Type  string
	Image string
	Flag  string
	Port  int32
}

var (
	instanceGVR = schema.GroupVersionResource{
		Group:    "isolet.dev",
		Version:  "v1",
		Resource: "instances",
	}
)

// WithTeamAndChallenge sets the team and challenge context for operations
func (s *instanceService) WithTeamAndChallenge(teamID, challengeID int64) *instanceService {
	s.teamID = teamID
	s.challengeID = challengeID
	return s
}

// WithChallenge sets challenge metadata for instance creation
func (s *instanceService) WithChallenge(challenge *ChallengeInfo) *instanceService {
	s.challenge = challenge
	return s
}

// WithInstance sets the instance object for update operations
func (s *instanceService) WithInstance(instance *instanceDom.Instance) *instanceService {
	s.instance = instance
	return s
}

// Start creates an Instance CRD for the Tide controller to manage
func (s *instanceService) Start(ctx context.Context) error {
	if s.teamID == 0 || s.challengeID == 0 {
		return fmt.Errorf("team ID and challenge ID must be set")
	}
	if s.challenge == nil {
		return fmt.Errorf("challenge info must be set")
	}

	// Generate instance name
	instanceName := generateInstanceName(s.teamID, s.challengeID)

	// Calculate expiry time
	expiresAt := time.Now().Add(s.config.Instances.Lifetime)

	// Build Instance CRD
	instance := s.buildInstanceCRD(instanceName, expiresAt)

	// Create the Instance CRD
	_, err := s.dynamicClient.Resource(instanceGVR).
		Namespace(s.namespace).
		Create(ctx, instance, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Instance CRD: %w", err)
	}

	return nil
}

// Stop deletes the Instance CRD, triggering cleanup by Tide controller
func (s *instanceService) Stop(ctx context.Context) error {
	if s.teamID == 0 || s.challengeID == 0 {
		return fmt.Errorf("team ID and challenge ID must be set")
	}

	instanceName := generateInstanceName(s.teamID, s.challengeID)

	err := s.dynamicClient.Resource(instanceGVR).
		Namespace(s.namespace).
		Delete(ctx, instanceName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete Instance CRD: %w", err)
	}

	return nil
}

// Extend updates the expiry time of an existing Instance CRD
func (s *instanceService) Extend(ctx context.Context) error {
	if s.teamID == 0 || s.challengeID == 0 {
		return fmt.Errorf("team ID and challenge ID must be set")
	}
	if s.instance == nil {
		return fmt.Errorf("instance must be set")
	}

	instanceName := generateInstanceName(s.teamID, s.challengeID)

	// Get current instance
	result, err := s.dynamicClient.Resource(instanceGVR).
		Namespace(s.namespace).
		Get(ctx, instanceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get Instance CRD: %w", err)
	}

	// Calculate new expiry (extend from current expiry, not from now)
	currentExpiry := time.Unix(s.instance.ExpiresAt, 0)
	newExpiry := currentExpiry.Add(s.config.Instances.Lifetime)
	
	// Check if exceeds max instance time
	if newExpiry.Sub(time.Now()) > s.config.Instances.MaxLifetime {
		return fmt.Errorf("extension would exceed maximum instance lifetime of %v", s.config.Instances.MaxLifetime)
	}

	// Update the expiry time
	err = unstructured.SetNestedField(result.Object, newExpiry.Format(time.RFC3339), "spec", "lifecycle", "expiresAt")
	if err != nil {
		return fmt.Errorf("failed to set new expiry time: %w", err)
	}

	// Update the Instance CRD
	_, err = s.dynamicClient.Resource(instanceGVR).
		Namespace(s.namespace).
		Update(ctx, result, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update Instance CRD: %w", err)
	}

	return nil
}

// GetStatus retrieves the current status of an Instance CRD
func (s *instanceService) GetStatus(ctx context.Context, teamID, challengeID int64) (*InstanceStatus, error) {
	instanceName := generateInstanceName(teamID, challengeID)

	result, err := s.dynamicClient.Resource(instanceGVR).
		Namespace(s.namespace).
		Get(ctx, instanceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Instance CRD: %w", err)
	}

	status := &InstanceStatus{
		Name:  instanceName,
		Phase: "Unknown",
	}

	// Extract status
	if statusMap, found, _ := unstructured.NestedMap(result.Object, "status"); found {
		if phase, ok := statusMap["phase"].(string); ok {
			status.Phase = phase
		}

		// Extract endpoints
		if endpoints, found, _ := unstructured.NestedSlice(result.Object, "status", "endpoints"); found {
			for _, ep := range endpoints {
				if epMap, ok := ep.(map[string]interface{}); ok {
					endpoint := Endpoint{
						Name:     getString(epMap, "name"),
						Protocol: getString(epMap, "protocol"),
					}
					if hostname, ok := epMap["hostname"].(string); ok {
						endpoint.Hostname = &hostname
					}
					if port, ok := epMap["port"].(int64); ok {
						port32 := int32(port)
						endpoint.Port = &port32
					}
					if ready, ok := epMap["ready"].(bool); ok {
						endpoint.Ready = ready
					}
					status.Endpoints = append(status.Endpoints, endpoint)
				}
			}
		}
	}

	return status, nil
}

// buildInstanceCRD constructs an Instance CRD spec
func (s *instanceService) buildInstanceCRD(name string, expiresAt time.Time) *unstructured.Unstructured {
	instance := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "isolet.dev/v1",
			"kind":       "Instance",
			"metadata": map[string]interface{}{
				"name": name,
				"labels": map[string]interface{}{
					"isolet.dev/team-id":         fmt.Sprintf("%d", s.teamID),
					"isolet.dev/challenge-id":    fmt.Sprintf("%d", s.challengeID),
					"isolet.dev/challenge-type":  s.challenge.Type,
				},
			},
			"spec": map[string]interface{}{
				"challenge": map[string]interface{}{
					"id":    fmt.Sprintf("%d", s.challenge.ID),
					"name":  s.challenge.Name,
					"type":  s.challenge.Type,
					"image": s.challenge.Image,
				},
				"team": map[string]interface{}{
					"id": s.teamID,
				},
				"endpoints": []interface{}{
					map[string]interface{}{
						"name":       "app",
						"protocol":   mapPortToProtocol(s.challenge.Port),
						"targetPort": s.challenge.Port,
					},
				},
				"lifecycle": map[string]interface{}{
					"expiresAt":      expiresAt.Format(time.RFC3339),
					"allowExtension": true,
				},
			},
		},
	}

	// Add flag if provided
	if s.challenge.Flag != "" {
		challenge, _, _ := unstructured.NestedMap(instance.Object, "spec", "challenge")
		challenge["flag"] = s.challenge.Flag
		unstructured.SetNestedMap(instance.Object, challenge, "spec", "challenge")
	}

	// Add resource requests/limits from config
	if s.config.Instances.DefaultCPU != "" || s.config.Instances.DefaultMemory != "" {
		requests := make(map[string]interface{})
		if s.config.Instances.DefaultCPU != "" {
			requests[string(corev1.ResourceCPU)] = s.config.Instances.DefaultCPU
		}
		if s.config.Instances.DefaultMemory != "" {
			requests[string(corev1.ResourceMemory)] = s.config.Instances.DefaultMemory
		}
		unstructured.SetNestedMap(instance.Object, requests, "spec", "requests")
	}

	if s.config.Instances.LimitCPU != "" || s.config.Instances.LimitMemory != "" {
		limits := make(map[string]interface{})
		if s.config.Instances.LimitCPU != "" {
			limits[string(corev1.ResourceCPU)] = s.config.Instances.LimitCPU
		}
		if s.config.Instances.LimitMemory != "" {
			limits[string(corev1.ResourceMemory)] = s.config.Instances.LimitMemory
		}
		unstructured.SetNestedMap(instance.Object, limits, "spec", "limits")
	}

	return instance
}

// Helper functions

func generateInstanceName(teamID, challengeID int64) string {
	return fmt.Sprintf("inst-team-%d-chall-%d", teamID, challengeID)
}

func mapPortToProtocol(port int32) string {
	switch port {
	case 80, 8080, 3000, 5000:
		return "http"
	case 443, 8443:
		return "https"
	case 22:
		return "ssh"
	default:
		return "nc"
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// InstanceStatus represents the status of an instance
type InstanceStatus struct {
	Name      string
	Phase     string
	Endpoints []Endpoint
}

// Endpoint represents a network endpoint
type Endpoint struct {
	Name     string
	Protocol string
	Hostname *string
	Port     *int32
	Ready    bool
}

func New() instanceDom.Service {
	config := utils.GetConfig()

	// Create Kubernetes config
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Errorf("failed to create in-cluster config: %w", err))
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(k8sConfig)
	if err != nil {
		panic(fmt.Errorf("failed to create dynamic client: %w", err))
	}

	return &instanceService{
		dynamicClient: dynamicClient,
		namespace:     config.Instances.Namespace,
		config:        config,
	}
}
