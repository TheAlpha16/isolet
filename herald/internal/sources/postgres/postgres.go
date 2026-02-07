package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/TheAlpha16/isolet/herald/internal/sources"
	"github.com/TheAlpha16/isolet/herald/pkg/facts"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/errors"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	slotName        = "herald_slot"
	publicationName = "herald_pub"
	pluginName      = "pgoutput"
)

type pgSource struct {
	sqlConn *pgx.Conn
	repl    *pgconn.PgConn

	lastLSN pglogrepl.LSN
	wg      *sync.WaitGroup
}

func (s *pgSource) Name() string {
	return utils.SourcePostgres
}

func (s *pgSource) Run(ctx context.Context, out chan<- facts.Fact) error {
	cfg := utils.GetConfig().Postgres

	conn, err := pgx.Connect(ctx, cfg.DSN)
	if err != nil {
		return errors.Raise(errors.ErrPostgresConnectionFailed, "", err)
	}
	s.sqlConn = conn

	replConn, err := pgconn.Connect(ctx, cfg.ReplicationDSN)
	if err != nil {
		return errors.Raise(errors.ErrPostgresConnectionFailed, "", err)
	}
	s.repl = replConn

	if err := s.ensurePublication(ctx); err != nil {
		return err
	}

	if err := s.ensureReplicationSlot(ctx); err != nil {
		return err
	}

	if err := pglogrepl.StartReplication(ctx, s.repl, slotName, s.lastLSN, pglogrepl.StartReplicationOptions{
		PluginArgs: []string{
			"proto_version '1'",
			fmt.Sprintf("publication_names '%s'", publicationName),
		},
	}); err != nil {
		return errors.Raise(errors.ErrPostgresReplicationFailed, "failed to start replication", err)
	}

	return nil
}

func NewSource(ctx context.Context, wg *sync.WaitGroup) sources.Source {
	return &pgSource{
		wg: wg,
	}
}
