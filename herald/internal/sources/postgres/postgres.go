package postgres

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/herald/internal/sources"
	"github.com/TheAlpha16/isolet/herald/pkg/facts"
	"github.com/TheAlpha16/isolet/herald/utils"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	slotName        = "herald_slot"
	publicationName = "herald_pub"
)

type pgSource struct {
	sqlPool *pgxpool.Pool
	repl    *pgconn.PgConn
	lastLSN pglogrepl.LSN

	wg *sync.WaitGroup
}

func (s *pgSource) Name() string {
	return utils.SourcePostgres
}

func (s *pgSource) Run(ctx context.Context, out chan<- facts.Fact) error {
	return nil
}

func NewSource(ctx context.Context, wg *sync.WaitGroup) sources.Source {
	return &pgSource{
		wg: wg,
	}
}
