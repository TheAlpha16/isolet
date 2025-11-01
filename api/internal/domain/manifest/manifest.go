package manifest

import (
	"github.com/TheAlpha16/isolet/api/internal/domain"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
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
	ID          int64
	ChallengeID int64
	Slug        string
	Image       string
	Type        challengeDom.ChallengeType
	Flag        *string
	Requests    []*Resource
	Limits      []*Resource
	Endpoints   []*Endpoint
	domain.BaseEntity
}

type Resource struct {
	ID         int64
	Name       ResourceName
	Value      string
	Type       ResourceType
	ManifestID int64
	domain.BaseEntity
}

type Endpoint struct {
	ID         int64
	ManifestID int64
	Name       string
	Protocol   Protocol
	TargetPort int32
	Hostname   *string
	Port       *int32
	domain.BaseEntity
}

func (m *Manifest) Populate(challenge *challengeDom.Challenge) {
	// TODO fininsh this
}
