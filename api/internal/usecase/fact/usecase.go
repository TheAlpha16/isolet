package fact

import (
	"context"

	factDom "github.com/TheAlpha16/isolet/api/internal/domain/fact"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
)

type factImpl struct {
	instanceUc instanceDom.Usecase
}

func (f *factImpl) HandleEvent(ctx context.Context, fact factDom.Fact) error {
	switch fact.FactType() {
	case factDom.InstanceExpiredFactType:
		if instFact, ok := fact.(*factDom.InstanceFact); ok {
			return f.handleInstanceExpired(ctx, instFact)
		}
	}
	return nil
}

func New(instanceUc instanceDom.Usecase) factDom.Usecase {
	return &factImpl{
		instanceUc: instanceUc,
	}
}
