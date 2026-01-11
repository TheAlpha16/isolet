package instance

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/api/internal/domain/manifest"
	"github.com/TheAlpha16/isolet/api/utils"
)

type instanceImpl struct {
	repo        instanceDom.Repository
	service     instanceDom.Service
	cache       cache.Cache
	challengeUc challengeDom.Usecase
	manifestUc  manifestDom.Usecase
	cvUc        cvDom.Usecase
}

func (i *instanceImpl) List(ctx context.Context) ([]*instanceDom.InstanceDTO, error) {
	// TODO finish implementation
	return nil, nil
}

func (i *instanceImpl) Start(ctx context.Context, input *instanceDom.StartInput) (*instanceDom.InstanceDTO, error) {
	config := utils.GetConfig()
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	// check if the challenge exists for this team and is of type on-demand
	challenge, err := i.challengeUc.ValidateAccess(ctx, input.ChallengeID, teamID)
	if err != nil {
		return nil, err
	}
	if challenge.Type != challengeDom.ChallengeOnDemand {
		return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotOnDemand, "", nil, nil)
	}

	// check if the instance is already running for the team and challenge
	instance, err := i.repo.GetByTeamAndChallenge(ctx, teamID, input.ChallengeID)
	if err != nil {
		if !errorDom.IsSameError(err, errorDom.ErrInstanceNotFound) {
			return nil, err
		}
	}
	if instance != nil {
		return instance.ToDTO(), nil
	}

	// take mutex on the instance to prevent race conditions
	if err := i.acquireInstanceLock(ctx, teamID, input.ChallengeID, config.Instances.StartTimeout); err != nil {
		return nil, err
	}
	defer i.releaseInstanceLock(ctx, teamID, input.ChallengeID) //nolint:errcheck

	manifest, err := i.manifestUc.GetByChallengeID(ctx, input.ChallengeID)
	if err != nil {
		return nil, err
	}

	if err := manifest.Populate(ctx, challenge); err != nil {
		return nil, err
	}

	var flag *string
	if manifest.FlagTemplate != nil {
		flag = utils.StringOrNil(challenge.RandomizedFlag(*manifest.FlagTemplate))
	}

	timeNow := time.Now()

	inst := instanceDom.Instance{
		TeamID:      utils.Int64OrNil(teamID),
		ChallengeID: manifest.ChallengeID,
		Flag:        flag,
		Lifecycle: &instanceDom.Lifecycle{
			ExpiresAt:      utils.TimePtr(timeNow.Add(utils.GetConfig().Instances.Lifetime)),
			AllowExtension: true,
		},
		Endpoints: []*instanceDom.Endpoint{}, // will be populated by service after K8s creation
	}

	if err := inst.Validate(ctx, manifest.Type); err != nil {
		return nil, err
	}

	if err := i.service.Start(ctx, &inst, manifest); err != nil {
		return nil, err
	}

	instance, err = i.repo.Create(ctx, &inst)
	if err != nil {
		return nil, err
	}

	return instance.ToDTO(), nil
}

func (i *instanceImpl) Stop(ctx context.Context, input *instanceDom.StopInput) error {
	config := utils.GetConfig()
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	// check if the instance exists for the team and challenge
	instance, err := i.repo.GetByID(ctx, input.InstanceID)
	if err != nil {
		return err
	}
	if instance.TeamID == nil || *instance.TeamID != teamID {
		return errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", nil, nil)
	}

	// take mutex on the instance to prevent race conditions
	if err := i.acquireInstanceLock(ctx, teamID, instance.ChallengeID, config.Instances.StopTimeout); err != nil {
		return err
	}
	defer i.releaseInstanceLock(ctx, teamID, instance.ChallengeID) //nolint:errcheck

	if err := i.service.Stop(ctx, instance); err != nil {
		return err
	}

	// delete from database
	if err := i.repo.Delete(ctx, instance.ID); err != nil {
		return err
	}

	return nil
}

func (i *instanceImpl) Extend(ctx context.Context, input *instanceDom.ExtendInput) (*instanceDom.InstanceDTO, error) {
	config := utils.GetConfig()
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	// check if the instance exists for the team and challenge
	instance, err := i.repo.GetByID(ctx, input.InstanceID)
	if err != nil {
		return nil, err
	}
	if instance.TeamID == nil || *instance.TeamID != teamID {
		return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", nil, nil)
	}

	// Get challenge to validate type
	challenge, err := i.challengeUc.ValidateAccess(ctx, instance.ChallengeID, teamID)
	if err != nil {
		return nil, err
	}

	// take mutex on the instance to prevent race conditions
	if err := i.acquireInstanceLock(ctx, teamID, instance.ChallengeID, config.Instances.ExtendTimeout); err != nil {
		return nil, err
	}
	defer i.releaseInstanceLock(ctx, teamID, instance.ChallengeID) //nolint:errcheck

	if err := instance.Validate(ctx, challenge.Type); err != nil {
		return nil, err
	}

	if !instance.Lifecycle.AllowExtension {
		return nil, errorDom.Raise(ctx, errorDom.ErrInstanceExtensionNotAllowed, "", nil, nil)
	}

	if instance.Lifecycle.ExpiresAt == nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrInstanceInvalid, "instance expiration time is nil", nil, common.ExtraData{"instance_id": instance.ID})
	}

	// check if the new expiry exceeds the maximum allowed duration
	newExpiry := instance.Lifecycle.ExpiresAt.Add(utils.GetConfig().Instances.Lifetime)

	if instance.CreatedAt.Add(utils.GetConfig().Instances.MaxLifetime).Before(newExpiry) {
		return nil, errorDom.Raise(ctx, errorDom.ErrInstanceExtensionNotAllowed, "instance exceeded maximum allowed duration", nil, common.ExtraData{"instance_id": instance.ID})
	}
	instance.Lifecycle.ExpiresAt = utils.TimePtr(newExpiry)

	if err := i.service.Extend(ctx, instance); err != nil {
		return nil, err
	}

	if err := i.repo.Update(ctx, instance, []string{"expires_at"}); err != nil {
		return nil, err
	}

	return instance.ToDTO(), nil
}

func (i *instanceImpl) acquireInstanceLock(ctx context.Context, teamID, challengeID int64, ttl time.Duration) error {
	key := instanceDom.InstanceCacheKey(teamID, challengeID)

	ok, err := i.cache.SetNXWithTTL(ctx, key, "1", ttl)
	if err != nil {
		return err
	}

	// if the key already exists, instance is already being started/stopped/updated
	if !ok {
		return errorDom.Raise(ctx, errorDom.ErrInstanceWIP, "", nil, common.ExtraData{"team_id": teamID, "challenge_id": challengeID})
	}

	return nil
}

func (i *instanceImpl) releaseInstanceLock(ctx context.Context, teamID, challengeID int64) error {
	return i.cache.Delete(ctx, instanceDom.InstanceCacheKey(teamID, challengeID))
}

func New(repo instanceDom.Repository, service instanceDom.Service, cache cache.Cache, challengeUc challengeDom.Usecase, manifestUc manifestDom.Usecase, cvUc cvDom.Usecase) instanceDom.Usecase {
	return &instanceImpl{
		repo:        repo,
		service:     service,
		cache:       cache,
		challengeUc: challengeUc,
		manifestUc:  manifestUc,
		cvUc:        cvUc,
	}
}
