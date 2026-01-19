package manifest

import (
	"context"

	"github.com/TheAlpha16/isolet/oracle/internal/domain"
	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/oracle/utils"
)

type ResourceName string
type Protocol string
type ResourceType string

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
	ProtocolNC    Protocol = "nc"
	ProtocolSSH   Protocol = "ssh"
)

const (
	ResourceCPU              ResourceName = "cpu"
	ResourceMemory           ResourceName = "memory"
	ResourceStorage          ResourceName = "storage"
	ResourceEphemeralStorage ResourceName = "ephemeral-storage"
)

const (
	ResourceTypeRequest ResourceType = "request"
	ResourceTypeLimit   ResourceType = "limit"
)

type Manifest struct {
	ID            int64
	ChallengeID   int64
	Slug          string
	Image         string
	Type          challengeDom.ChallengeType
	FlagTemplate  *string
	Requests      []*Resource
	Limits        []*Resource
	EndpointSpecs []*EndpointSpec
	domain.BaseEntity
}

func (m *Manifest) Populate(ctx context.Context, challenge *challengeDom.Challenge) error {
	slug, err := challenge.Slug(ctx)
	if err != nil {
		return err
	}

	m.Slug = slug
	m.Type = challenge.Type

	// apply default resource limits if not set
	limits := map[ResourceName]string{}

	for _, limit := range m.Limits {
		limits[limit.Name] = limit.Value
	}

	if _, ok := limits[ResourceCPU]; !ok {
		m.Limits = append(m.Limits, &Resource{
			Name:  ResourceCPU,
			Type:  ResourceTypeLimit,
			Value: utils.GetConfig().Instances.LimitCPU,
		})
	}

	if _, ok := limits[ResourceMemory]; !ok {
		m.Limits = append(m.Limits, &Resource{
			Name:  ResourceMemory,
			Type:  ResourceTypeLimit,
			Value: utils.GetConfig().Instances.LimitMemory,
		})
	}

	return nil
}

type Resource struct {
	ID         int64
	Name       ResourceName
	Value      string
	Type       ResourceType
	ManifestID int64
	domain.BaseEntity
}

type EndpointSpec struct {
	ID         int64
	Name       string
	Protocol   Protocol
	TargetPort int32
	ManifestID int64
	domain.BaseEntity
}
