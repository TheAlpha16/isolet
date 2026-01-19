package event

import (
	"context"

	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	eventDom "github.com/TheAlpha16/isolet/oracle/internal/domain/event"
)

type eventImpl struct {
	cvUc cvDom.Usecase
}

func (e *eventImpl) Info(ctx context.Context) (*eventDom.InfoOutput, error) {
	return &eventDom.InfoOutput{
		Name:       e.cvUc.GetString(ctx, cvDom.EventName),
		StartTime:  int64(e.cvUc.GetInt(ctx, cvDom.EventStart)),
		EndTime:    int64(e.cvUc.GetInt(ctx, cvDom.EventEnd)),
		PostEvent:  e.cvUc.GetBool(ctx, cvDom.EventPostMode),
		TeamLength: e.cvUc.GetInt(ctx, cvDom.TeamMaxSize),
	}, nil
}

func New(cvUc cvDom.Usecase) eventDom.Usecase {
	return &eventImpl{
		cvUc: cvUc,
	}
}
