package instance

import (
	"fmt"

	"github.com/TheAlpha16/isolet/api/internal/domain"
)

type Protocol string

const (
	InstanceCachePrefix = "instance"
)

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
	ProtocolNC    Protocol = "nc"
	ProtocolSSH   Protocol = "ssh"
)

type Instance struct {
	ID          int64
	ChallengeID int64
	TeamID      int64
	ExpiresAt   int64
	Info        *DeploymentInfo
	domain.BaseEntity
}

type DeploymentInfo struct {
	ID         int64
	Slug       string
	Image      string
	TargetPort int32
	Requests   []*Resource
	Limits     []*Resource
	Endpoints  []*Endpoint
}

type Resource struct {
	Name  string
	Value string
}

type Endpoint struct {
	ID         int64
	Name       string
	Protocol   Protocol
	TargetPort int32
	Hostname   *string
	Port       *int32
}

func InstanceCacheKey(teamID, challengeID int64) string {
	return fmt.Sprintf("%s:team:%d:challenge:%d", InstanceCachePrefix, teamID, challengeID)
}
