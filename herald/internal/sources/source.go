package sources

import (
	"context"

	"github.com/TheAlpha16/isolet/herald/pkg/facts"
)

type Source interface {
	Name() string
	Run(ctx context.Context, out chan<- facts.Fact) error
}
