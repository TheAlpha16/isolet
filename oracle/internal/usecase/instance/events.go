package instance

import (
	"context"

	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"

	tideConstants "github.com/TheAlpha16/isolet/tide/utils"
)

func (i *instanceImpl) HandleEvent(ctx context.Context, eventType string, instance *instanceDom.Instance) error {
	var handlerMap = map[string]func(context.Context, *instanceDom.Instance) error{
		tideConstants.EventReasonExpired: i.handleInstanceExpired,
	}

	handler, exists := handlerMap[eventType]
	if !exists {
		return nil
	}

	return handler(ctx, instance)
}

func (i instanceImpl) handleInstanceExpired(ctx context.Context, instance *instanceDom.Instance) error {
	return i.repo.DeleteByRefs(ctx, instance.TeamID, instance.ChallengeID)
}
