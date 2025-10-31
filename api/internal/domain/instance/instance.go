package instance

import (
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
)

const (
	InstanceCachePrefix = "instance"
)

type Instance struct {
	ID        int64
	TeamID    int64
	Manifest  *manifestDom.Manifest
	Lifecycle *Lifecycle
	domain.BaseEntity
}

type Lifecycle struct {
	AvailableAt    *time.Time
	ExpiresAt      *time.Time
	AllowExtension bool
}

func InstanceCacheKey(teamID, challengeID int64) string {
	return fmt.Sprintf("%s:team:%d:challenge:%d", InstanceCachePrefix, teamID, challengeID)
}
