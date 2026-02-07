package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/TheAlpha16/isolet/herald/internal/sources"
	"github.com/TheAlpha16/isolet/herald/pkg/facts"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/tracer"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

const (
	slotName        = "herald_slot"
	publicationName = "herald_pub"
	pluginName      = "pgoutput"
)

var pgTracer = otel.Tracer("herald.sources.postgres")

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
	ctx, span, log := tracer.StartSpan(ctx, pgTracer, "herald.sources.postgres.Run")
	defer span.End()

	if err := s.ensurePublication(ctx); err != nil {
		errors.HandleSpanError(ctx, span, log, "failed to ensure publication", err)
		return err
	}

	if err := s.ensureReplicationSlot(ctx); err != nil {
		errors.HandleSpanError(ctx, span, log, "failed to ensure replication slot", err)
		return err
	}

	if err := pglogrepl.StartReplication(ctx, s.repl, slotName, s.lastLSN, pglogrepl.StartReplicationOptions{
		PluginArgs: []string{
			"proto_version '1'",
			fmt.Sprintf("publication_names '%s'", publicationName),
		},
	}); err != nil {
		err = errors.Raise(errors.ErrPostgresReplicationFailed, "failed to start replication", err)
		errors.HandleSpanError(ctx, span, log, "failed to start replication", err)
		return err
	}

	return nil
}

func NewSource(ctx context.Context, wg *sync.WaitGroup) sources.Source {
	ctx, span, log := tracer.StartSpan(ctx, pgTracer, "herald.sources.postgres.NewSource")
	defer span.End()

	cfg := utils.GetConfig().Postgres

	conn, err := pgx.Connect(ctx, cfg.DSN)
	if err != nil {
		err = errors.Raise(errors.ErrPostgresConnectionFailed, "", err)
		errors.HandleSpanError(ctx, span, log, "failed to connect to postgres", err)
		log.Fatal("failed to connect to postgres", zap.Error(err))
	}

	replConn, err := pgconn.Connect(ctx, cfg.ReplicationDSN)
	if err != nil {
		err = errors.Raise(errors.ErrPostgresConnectionFailed, "", err)
		errors.HandleSpanError(ctx, span, log, "failed to connect to postgres replication", err)
		log.Fatal("failed to connect to postgres replication", zap.Error(err))
	}

	return &pgSource{
		sqlConn: conn,
		repl:    replConn,
		wg:      wg,
	}
}
