package emitter

import (
	"context"

	"github.com/TheAlpha16/isolet/herald/pkg/facts"
)

type Emitter interface {
	Emit(ctx context.Context, fact facts.Fact) error
	Close() error
}
