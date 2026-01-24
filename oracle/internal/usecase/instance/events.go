package instance

import (
	"context"
	"fmt"

	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	"github.com/TheAlpha16/isolet/oracle/utils/tracer"

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
	deletedRows, err := i.repo.DeleteByRefs(ctx, instance.TeamID, instance.ChallengeID)
	if err != nil {
		return err
	}

	if deletedRows > 0 {
		tracer.InstanceActiveTotal.WithLabelValues(fmt.Sprintf("%d", instance.ChallengeID)).Dec()
	}

	return nil
}
