package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TheAlpha16/isolet/herald/pkg/facts"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/tracer"
	"github.com/google/uuid"

	"go.uber.org/zap"
)

type InstanceRow struct {
	ID          int64  `json:"id"`
	ChallengeID int64  `json:"challenge_id"`
	TeamID      *int64 `json:"team_id"`
	ExpiresAt   int64  `json:"expires_at"`
}

func instancesHandler(data map[string]any, factChan chan<- facts.Fact, eventType eventType) {
	ctx, span, log := tracer.StartSpan(context.Background(), pgTracer, "herald.sources.postgres.instancesHandler")
	defer span.End()

	log = log.With(
		zap.String("event_type", string(eventType)),
		zap.String("table", string(tableInstances)),
	)

	instance, err := parseInstanceRow(data)
	if err != nil {
		errors.HandleSpanError(ctx, span, log, "failed to parse instance row", err)
		return
	}

	log.Debug(
		"received instance event",
		zap.Int64("instance_id", instance.ID),
		zap.Int64("challenge_id", instance.ChallengeID),
		zap.Int64p("team_id", instance.TeamID),
		zap.Int64("expires_at", instance.ExpiresAt),
	)

	fact := &facts.Notification{
		ID: uuid.NewString(),
		BaseFact: facts.BaseFact{
			Type: facts.NotificationFactType,
			At:   time.Now(),
		},
		Entity: &facts.Entity{
			Name: facts.EntityInstance,
			ID:   instance.ID,
			Data: func() json.RawMessage {
				b, err := json.Marshal(instance)
				if err != nil {
					errors.HandleSpanError(ctx, span, log, "failed to marshal instance data", err)
					return nil
				}
				return b
			}(),
		},
		TeamIDs: []int64{},
	}

	var action facts.Action
	if a, ok := actionMap[eventType]; ok {
		action = a
	}
	fact.Action = &action

	if instance.TeamID != nil {
		fact.TeamIDs = append(fact.TeamIDs, *instance.TeamID)
	}

	if err := fact.Validate(); err != nil {
		errors.HandleSpanError(ctx, span, log, "failed to validate NotificationFact", err)
		return
	}

	factChan <- fact
	span.AddEvent("postgres.fact.buffered")
	log.Debug(
		"sent fact to channel",
		zap.ByteString("key", fact.Key()),
		zap.String("fact_type", string(fact.FactType())),
	)
}
