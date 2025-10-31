package instance

import (
	"context"
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain"
	"github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
)

const (
	InstanceCachePrefix = "instance"
)

type Instance struct {
	ID        int64
	TeamID    *int64
	Manifest  *manifestDom.Manifest
	Lifecycle *Lifecycle
	domain.BaseEntity
}

func (in *Instance) Validate(ctx context.Context) error {
	extraData := common.ExtraData{"instance": *in}

	// manifest
	if in.Manifest == nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "instance manifest cannot be nil", nil, extraData)
	}

	// lifecycle
	if in.Lifecycle == nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "instance lifecycle cannot be nil", nil, extraData)
	}

	// challenge
	if in.Manifest.ChallengeID == 0 {
		return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "challenge ID in manifest cannot be zero", nil, extraData)
	}

	switch in.Manifest.Type {
	case challenge.ChallengeStatic:
		return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "instance cannot be spawned for static challenges", nil, extraData)

	case challenge.ChallengeOnDemand:
		if in.TeamID == nil {
			return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "team ID cannot be nil for on-demand challenge instances", nil, extraData)
		}
		if in.Lifecycle.ExpiresAt == nil {
			return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "lifecycle expires_at cannot be nil for on-demand challenge instances", nil, extraData)
		}
	case challenge.ChallengeDynamic:
		if in.TeamID != nil {
			return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "team ID must be nil for dynamic challenge instances", nil, extraData)
		}
	default:
		return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "invalid challenge type in manifest", nil, extraData)
	}

	return nil
}

type Lifecycle struct {
	AvailableAt    *time.Time
	ExpiresAt      *time.Time
	AllowExtension bool
}

func InstanceCacheKey(teamID, challengeID int64) string {
	return fmt.Sprintf("%s:team:%d:challenge:%d", InstanceCachePrefix, teamID, challengeID)
}
