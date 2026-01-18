package fact

import (
	"context"
	"time"
)

type FactType string

type Fact interface {
	FactType() FactType
	Key() []byte
	OccuredAt() time.Time
	Validate() error
}

type Usecase interface {
	HandleEvent(ctx context.Context, fact Fact) error
}
