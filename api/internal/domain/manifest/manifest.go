package manifest

import (
	"context"

	"github.com/TheAlpha16/isolet/api/internal/domain"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/api/utils"
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

func (m *Manifest) Populate(ctx context.Context, challenge *challengeDom.Challenge) error {
	slug, err := challenge.Slug(ctx)
	if err != nil {
		return err
	}

	m.Slug = slug
	m.Type = challenge.Type

	if m.Flag != nil {
		m.Flag = utils.StringOrNil(challenge.RandomizedFlag(*m.Flag))
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
