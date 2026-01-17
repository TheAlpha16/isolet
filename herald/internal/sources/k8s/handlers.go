package k8s

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/utils/logger"
	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"github.com/google/uuid"
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

	fact := &facts.InstanceFact{
		BaseFact: facts.BaseFact{
			Type: facts.InstanceExpiredFactType,
			At:   time.Now(),
		},
		ID: uuid.New(),
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
