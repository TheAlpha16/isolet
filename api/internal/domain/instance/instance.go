package instance

import (
	"fmt"

	"github.com/TheAlpha16/isolet/api/internal/domain"
)

const (
	InstanceCachePrefix = "instance"
)

type Instance struct {
	ID          int64
	ChallengeID int64
	TeamID      int64
	ExpiresAt   int64
	domain.BaseEntity
}

func InstanceCacheKey(teamID, challengeID int64) string {
	return fmt.Sprintf("%s:team:%d:challenge:%d", InstanceCachePrefix, teamID, challengeID)
}
