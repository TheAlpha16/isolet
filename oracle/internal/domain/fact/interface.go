package fact

import (
	"context"
)

type Usecase interface {
	HandleEvent(ctx context.Context, fact Fact) error
}
