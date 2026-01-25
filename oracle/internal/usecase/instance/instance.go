package instance

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/TheAlpha16/isolet/oracle/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/oracle/internal/domain/manifest"
	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/tracer"

	"go.opentelemetry.io/otel"
)

var instanceTracer = otel.Tracer("usecase.instance")

type instanceImpl struct {
	repo    instanceDom.Repository
	service instanceDom.Service

	challengeUc challengeDom.Usecase
	manifestUc  manifestDom.Usecase
	cvUc        cvDom.Usecase

	cache cache.Cache
	wg    *sync.WaitGroup
}

func (i *instanceImpl) List(ctx context.Context) ([]*instanceDom.InstanceDTO, error) {
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	instances, err := i.repo.GetByTeam(ctx, teamID)
	if err != nil {
		return nil, err
	}

	var instanceDTOs []*instanceDom.InstanceDTO
	for _, inst := range instances {
		instanceDTOs = append(instanceDTOs, inst.ToDTO())
	}

	return instanceDTOs, nil
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
	instance, err := i.repo.GetByRefs(ctx, &teamID, input.ChallengeID)
	if err != nil {
		if !errorDom.Is(err, errorDom.ErrInstanceNotFound) {
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
		Domain:    i.cvUc.GetString(ctx, cvDom.InstanceDomain),
	}

	if err := inst.Validate(ctx, manifest.Type); err != nil {
		return nil, err
	}

	if err := i.service.Start(ctx, &inst, manifest); err != nil {
		tracer.InstanceOperationsTotal.WithLabelValues(
			strconv.FormatInt(input.ChallengeID, 10),
			tracer.OpStart,
			tracer.StatusFailure,
		).Inc()
		return nil, err
	}

	instance, err = i.repo.Create(ctx, &inst)
	if err != nil {
		tracer.InstanceOperationsTotal.WithLabelValues(
			strconv.FormatInt(input.ChallengeID, 10),
			tracer.OpStart,
			tracer.StatusFailure,
		).Inc()
		return nil, err
	}

	tracer.InstanceOperationsTotal.WithLabelValues(
		strconv.FormatInt(input.ChallengeID, 10),
		tracer.OpStart,
		tracer.StatusSuccess,
	).Inc()

	tracer.InstanceProvisionDurationSeconds.WithLabelValues(
		strconv.FormatInt(input.ChallengeID, 10),
		tracer.OpStart,
		tracer.StatusSuccess,
	).Observe(time.Since(timeNow).Seconds())

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
		// if instance doesn't exist in K8s, still proceed to clean up DB record
		if !errorDom.Is(err, errorDom.ErrK8sInstanceNotFound) {
			return err
		}
	}

	// delete from database
	if err := i.repo.Delete(ctx, instance.ID); err != nil {
		return err
	}

	tracer.InstanceOperationsTotal.WithLabelValues(
		strconv.FormatInt(instance.ChallengeID, 10),
		tracer.OpStop,
		tracer.StatusSuccess,
	).Inc()

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

	// update instance domain from config vars
	instance.Domain = i.cvUc.GetString(ctx, cvDom.InstanceDomain)

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
		tracer.InstanceOperationsTotal.WithLabelValues(
			strconv.FormatInt(instance.ChallengeID, 10),
			tracer.OpExtend,
			tracer.StatusFailure,
		).Inc()
		if errorDom.Is(err, errorDom.ErrK8sInstanceNotFound) {
			// delete the orphaned DB record
			if delErr := i.repo.Delete(ctx, instance.ID); delErr == nil {
				tracer.InstanceOperationsTotal.WithLabelValues(
					strconv.FormatInt(instance.ChallengeID, 10),
					tracer.OpStop,
					tracer.StatusSuccess,
				).Inc()
			}
			return nil, errorDom.Raise(ctx, errorDom.ErrInstanceNotFound, "", err, nil)
		}
		return nil, err
	}

	if err := i.repo.Update(ctx, instance.ID, map[string]any{
		instanceDom.ExpiresAtColumn: instance.Lifecycle.ExpiresAt.Unix(),
	}); err != nil {
		tracer.InstanceOperationsTotal.WithLabelValues(
			strconv.FormatInt(instance.ChallengeID, 10),
			tracer.OpExtend,
			tracer.StatusFailure,
		).Inc()
		return nil, err
	}

	tracer.InstanceOperationsTotal.WithLabelValues(
		strconv.FormatInt(instance.ChallengeID, 10),
		tracer.OpExtend,
		tracer.StatusSuccess,
	).Inc()

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

func (i *instanceImpl) runInstanceExpiryCleanup(ctx context.Context, interval time.Duration) {
	i.wg.Go(func() {
		ctx, span, logger := tracer.StartSpan(ctx, instanceTracer, "instanceImpl.InstanceCleanup")
		defer span.End()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				challengeCounts, err := i.repo.DeleteExpired(ctx, time.Now())
				if err != nil {
					errorDom.HandleSpanError(ctx, span, logger, "failed to delete expired instances", err)
					continue
				}

				for challengeID, count := range challengeCounts {
					tracer.InstanceOperationsTotal.WithLabelValues(
						strconv.FormatInt(challengeID, 10),
						tracer.OpExpire,
						tracer.StatusSuccess,
					).Add(float64(count))
				}

			case <-ctx.Done():
				return
			}
		}
	})
}

func New(ctx context.Context, repo instanceDom.Repository, service instanceDom.Service, cache cache.Cache, challengeUc challengeDom.Usecase, manifestUc manifestDom.Usecase, cvUc cvDom.Usecase, wg *sync.WaitGroup) instanceDom.Usecase {
	ctx, cancel := context.WithCancel(ctx)

	inst := &instanceImpl{
		repo:        repo,
		service:     service,
		cache:       cache,
		challengeUc: challengeUc,
		manifestUc:  manifestUc,
		cvUc:        cvUc,
		wg:          wg,
	}

	inst.runInstanceExpiryCleanup(ctx, utils.GetConfig().Instances.ExpiryCheckInterval)

	utils.InterruptHandlerChannel <- func() {
		cancel()
	}

	return inst
}
