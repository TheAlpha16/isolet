package instance

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"

	tidev1 "github.com/TheAlpha16/isolet/tide/api/v1"
	"go.uber.org/zap"
)

func (is *instanceSvc) ensureInstanceReady(ctx context.Context, old *tidev1.Instance, terminalStates ...tidev1.Phase) error {
	var createdInst *tidev1.Instance
	var err error

	stateMap := make(map[tidev1.Phase]struct{})
	for _, state := range terminalStates {
		stateMap[state] = struct{}{}
	}

	createdInst, err = is.client.GetInstance(ctx, old.Name, old.Namespace)
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceCreationFailed, "failed to get created instance", err, nil)
	}

	for {
		createdInst, err = is.client.GetInstance(ctx, old.Name, old.Namespace)
		if err != nil {
			return errorDom.Raise(ctx, errorDom.ErrInstanceCreationFailed, "failed to get created instance", err, nil)
		}
		if _, ok := stateMap[createdInst.Status.Phase]; ok {
			// process endpoints
			allReady := true
			for _, ep := range createdInst.Status.Endpoints {
				if !ep.Ready {
					logger.GetLogger(ctx).Info("endpoint not ready yet",
						zap.String("instance", old.Name),
						zap.String("endpoint", ep.Name),
					)
					allReady = false
					break
				}
			}
			if allReady {
				break
			}
		}
		logger.GetLogger(ctx).Info("waiting for instance to reach terminal state",
			zap.String("instance", old.Name),
			zap.String("phase", string(createdInst.Status.Phase)),
		)
		time.Sleep(utils.GetConfig().K8s.InstancePollRate)
	}

	if createdInst.Status.Phase != tidev1.PhaseRunning {
		return errorDom.Raise(ctx, errorDom.ErrInstanceCreationFailed, "instance reached failed state", nil, common.ExtraData{
			"instance_name": createdInst.Name,
			"phase":         createdInst.Status.Phase,
		})
	}

	old.Status = createdInst.Status

	return nil
}
