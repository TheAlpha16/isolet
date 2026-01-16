package instance

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"

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
	existingInstance, err := i.repo.GetByRefs(ctx, instance.TeamID, instance.ChallengeID)
	if err != nil {
		if errorDom.IsSameError(err, errorDom.ErrInstanceNotFound) {
			return nil
		}
		return err
	}

	if err := i.repo.Delete(ctx, existingInstance.ID); err != nil {
		return err
	}

	return nil
}
