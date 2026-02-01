package postgres

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/herald/internal/sources"
	"github.com/TheAlpha16/isolet/herald/pkg/facts"
	"github.com/TheAlpha16/isolet/herald/utils"
)

type pgSource struct {
}

func (s *pgSource) Name() string {
	return utils.SourcePostgres
}

func (s *pgSource) Run(ctx context.Context, out chan<- facts.Fact) error {
	return nil
}

func NewSource(ctx context.Context, wg *sync.WaitGroup) sources.Source {
	return &pgSource{}
}
