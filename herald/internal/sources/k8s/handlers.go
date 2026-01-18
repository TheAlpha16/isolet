package k8s

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/cache"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/tracer"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"go.uber.org/zap"
)

func handleInstance(obj any, out chan<- facts.Fact, eventType eventType) {
	ctx, span, log := tracer.StartSpan(context.Background(), k8sTracer, "herald.sources.k8s.handleInstance")
	defer span.End()

	if eventType != eventTypeUpdate {
		return // dont care about create/delete events for now
	}

	instance, ok := extractObject[*tidev1.Instance](obj)
	if !ok {
		log.Warn("failed to extract instance from object")
		return
	}

	phase := instance.Status.Phase

	log = log.With(
		zap.String("event_type", string(eventType)),
		zap.ByteString("uid", []byte(instance.UID)),
		zap.String("phase", string(phase)),
	)

	log.Debug("received instance event from k8s",
		zap.String("instance", instance.Name),
		zap.String("namespace", instance.Namespace),
		zap.String("resource_version", instance.ResourceVersion),
	)

	if phase != tidev1.PhaseExpired {
		return
	}

	cacheKey := getCacheKey(string(tidev1.PhaseExpired), string(instance.UID))
	success, err := cache.SetNXWithTTL(ctx, cacheKey, "1", utils.GetConfig().K8s.DeDupeWindow)
	if err != nil {
		errors.HandleSpanError(
			ctx, span, log,
			"failed to set de-dupe cache key", err,
			zap.String("cache_key", cacheKey),
		)
	} else if !success {
		log.Debug("received duplicate event")
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
		errors.HandleSpanError(
			ctx, span, log,
			"failed to validate InstanceFact", err,
		)
		return
	}

	out <- fact
	span.AddEvent("instance.fact.buffered")
	log.Debug(
		"sent fact to channel",
		zap.ByteString("key", fact.Key()),
		zap.String("fact_type", string(fact.FactType())),
	)
}
