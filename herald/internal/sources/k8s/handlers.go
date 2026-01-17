package k8s

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/cache"
	"github.com/TheAlpha16/isolet/herald/utils/logger"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
)

func handleInstance(obj any, out chan<- facts.Fact, eventType eventType) {
	if eventType != eventTypeUpdate {
		return // dont care about create/delete events for now
	}

	ctx, span := tracer.Start(context.Background(), "herald.sources.k8s.handleInstance")
	defer span.End()

	instance, ok := extractObject[*tidev1.Instance](obj)
	if !ok {
		return
	}

	phase := instance.Status.Phase

	if phase != tidev1.PhaseExpired {
		return
	}

	cacheKey := getCacheKey(string(tidev1.PhaseExpired), string(instance.UID))
	success, _ := cache.SetNXWithTTL(ctx, cacheKey, "1", utils.GetConfig().K8s.DeDupeWindow)
	if !success {
		// already processed recently
		return
	}

	fact := &facts.InstanceFact{
		BaseFact: facts.BaseFact{
			Type: facts.InstanceExpiredFactType,
			At:   time.Now(),
		},
		ID:          string(instance.UID),
		ChallengeID: instance.Spec.Challenge.ID,
	}

	if instance.Spec.Team != nil {
		var teamId int64 = instance.Spec.Team.ID
		fact.TeamID = &teamId
	}

	if err := fact.Validate(); err != nil {
		handleError(ctx, span, logger.GetAppLogger(), "failed to validate InstanceFact", err, map[string]any{
			"instance":  instance.Name,
			"namespace": instance.Namespace,
			"phase":     phase,
		})
		return
	}

	out <- fact
}
