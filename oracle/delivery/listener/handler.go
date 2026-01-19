package listener

import (
	"context"

	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	"github.com/TheAlpha16/isolet/oracle/internal/usecase"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"

	tideConstants "github.com/TheAlpha16/isolet/tide/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func handleDeleteEvent(obj interface{}, usecases *usecase.Usecases) {
	ctx, span := tracer.Start(context.Background(), "listener.handleDeleteEvent", trace.WithAttributes(
		attribute.String("event.type", tideConstants.EventReasonExpired),
	))
	defer span.End()

	logger := logger.GetAppLogger()

	instance, fromTombstone, err := toInstance(ctx, obj)
	if err != nil {
		handleError(ctx, span, logger, "failed to extract instance from event object", err, nil)
		return
	}

	if fromTombstone {
		span.SetAttributes(attribute.Bool("instance.from_tombstone", true))
	}

	span.SetAttributes(
		attribute.String("instance.name", instance.Name),
		attribute.String("instance.namespace", instance.Namespace),
		attribute.Int64("instance.challenge_id", instance.Spec.Challenge.ID),
	)

	instDom := &instanceDom.Instance{
		ChallengeID: instance.Spec.Challenge.ID,
	}
	if instance.Spec.Team != nil {
		instDom.TeamID = &instance.Spec.Team.ID
		span.SetAttributes(attribute.Int64("instance.team_id", instance.Spec.Team.ID))
	}

	span.AddEvent("instance.delete_event_received")

	if err := usecases.Instance.HandleEvent(ctx, tideConstants.EventReasonExpired, instDom); err != nil {
		handleError(ctx, span, logger, "failed to handle instance expiry event", err, map[string]any{
			"instance":     instance.Name,
			"challenge_id": instance.Spec.Challenge.ID,
		})
	} else {
		span.SetStatus(codes.Ok, "instance deletion event handled successfully")
		span.AddEvent("instance.delete_event_handled")
		logger.Info("Handled instance deletion event",
			zap.String("instance", instance.Name),
			zap.Int64("challenge_id", instance.Spec.Challenge.ID),
		)
	}
}
