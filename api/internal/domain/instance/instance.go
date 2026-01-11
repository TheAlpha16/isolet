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
	"github.com/TheAlpha16/isolet/api/utils"
)

const (
	InstanceCachePrefix = "instance"
)

type Instance struct {
	ID          int64
	TeamID      *int64
	ChallengeID int64
	Flag        *string
	Lifecycle   *Lifecycle
	Endpoints   []*Endpoint
	domain.BaseEntity
}

func (in *Instance) Validate(ctx context.Context, challengeType challenge.ChallengeType) error {
	extraData := common.ExtraData{"instance": *in}

	// lifecycle
	if in.Lifecycle == nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "instance lifecycle cannot be nil", nil, extraData)
	}

	// challenge
	if in.ChallengeID == 0 {
		return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "challenge ID cannot be zero", nil, extraData)
	}

	switch challengeType {
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
		return errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "invalid challenge type", nil, extraData)
	}

	return nil
}

func (in *Instance) Name() string {
	var identifier string
	if in.TeamID == nil {
		identifier = fmt.Sprintf("dynamic@%d", in.ChallengeID)
	} else {
		identifier = fmt.Sprintf("team%d@%d", *in.TeamID, in.ChallengeID)
	}
	return utils.HMAC256(identifier, utils.GetConfig().Instances.SecretKey)[:16]
}

type Lifecycle struct {
	AvailableAt    *time.Time
	ExpiresAt      *time.Time
	AllowExtension bool
}

type Endpoint struct {
	ID         int64
	InstanceID int64
	Name       string
	Protocol   manifestDom.Protocol
	TargetPort int32
	Hostname   *string
	Port       *int32
	Ready      bool
	domain.BaseEntity
}

func InstanceCacheKey(teamID, challengeID int64) string {
	return fmt.Sprintf("%s:team:%d:challenge:%d", InstanceCachePrefix, teamID, challengeID)
}
