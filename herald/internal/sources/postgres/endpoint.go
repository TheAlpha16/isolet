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

type EndpointRow struct {
	ID         int64  `json:"id"`
	InstanceID int64  `json:"instance_id"`
	TeamID     *int64 `json:"team_id,omitempty"`
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Hostname   string `json:"hostname"`
	Port       *int32 `json:"port,omitempty"`
	Ready      bool   `json:"ready"`
}

func endpointsHandler(data map[string]any, factChan chan<- facts.Fact, eventType eventType) {
	ctx, span, log := tracer.StartSpan(context.Background(), pgTracer, "herald.sources.postgres.endpointsHandler")
	defer span.End()

	log = log.With(
		zap.String("event_type", string(eventType)),
		zap.String("table", string(tableInstances)),
	)

	endpoint, err := parseEndpointRow(data)
	if err != nil {
		errors.HandleSpanError(ctx, span, log, "failed to parse endpoint row", err)
		return
	}

	log.Debug(
		"received endpoint event",
		zap.Int64("endpoint_id", endpoint.ID),
		zap.Int64("instance_id", endpoint.InstanceID),
		zap.Int64p("team_id", endpoint.TeamID),
	)

	fact := &facts.Notification{
		ID: uuid.NewString(),
		BaseFact: facts.BaseFact{
			Type: facts.NotificationFactType,
			At:   time.Now(),
		},
		Entity: &facts.Entity{
			Name: facts.EntityEndpoint,
			ID:   endpoint.ID,
			Data: func() json.RawMessage {
				b, err := json.Marshal(endpoint)
				if err != nil {
					errors.HandleSpanError(ctx, span, log, "failed to marshal endpoint data", err)
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

	if endpoint.TeamID != nil {
		fact.TeamIDs = append(fact.TeamIDs, *endpoint.TeamID)
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
