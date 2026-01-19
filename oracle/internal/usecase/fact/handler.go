package fact

import (
	"context"

	factDom "github.com/TheAlpha16/isolet/oracle/internal/domain/fact"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	tideConstants "github.com/TheAlpha16/isolet/tide/utils"
)

func (f *factImpl) handleInstanceExpired(ctx context.Context, fact *factDom.InstanceFact) error {
	inst := &instanceDom.Instance{
		TeamID:      fact.TeamID,
		ChallengeID: fact.ChallengeID,
	}

	return f.instanceUc.HandleEvent(ctx, tideConstants.EventReasonExpired, inst)
}
