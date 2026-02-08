package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TheAlpha16/isolet/herald/pkg/facts"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/tracer"
	"go.uber.org/zap"
)

type table string

const (
	tableInstances table = "instances"
)

var tableHandlerMap = map[table]func(data map[string]any, factChan chan<- facts.Fact, eventType eventType){
	tableInstances: instancesHandler,
}

var actionMap = map[eventType]facts.Action{
	eventTypeCreate: facts.ActionCreated,
	eventTypeUpdate: facts.ActionUpdated,
	eventTypeDelete: facts.ActionDeleted,
}

func instancesHandler(data map[string]any, factChan chan<- facts.Fact, eventType eventType) {
	ctx, span, log := tracer.StartSpan(context.Background(), pgTracer, "herald.sources.postgres.instancesHandler")
	defer span.End()

	log = log.With(
		zap.String("event_type", string(eventType)),
		zap.String("table", string(tableInstances)),
	)

	// DEBUG
	log.Debug("handling postgres event", zap.Any("data", data))

	fact := &facts.Notification{
		BaseFact: facts.BaseFact{
			Type: facts.NotificationFactType,
			At:   time.Now(),
		},
		Entity: &facts.Entity{
			Name: facts.EntityInstance,
			ID:   data["id"].(int64),
			Data: json.RawMessage{},
		},
		TeamIDs: []int64{
			data["team_id"].(int64),
		},
	}

	var action facts.Action
	if a, ok := actionMap[eventType]; ok {
		action = a
	}
	fact.Action = &action

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
