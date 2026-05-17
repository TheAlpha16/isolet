package postgres

import (
	"context"
	goerrors "errors"
	"fmt"
	"strings"

	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (s *pgSource) ensurePublication(ctx context.Context) error {
	var exists bool

	if err := s.sqlConn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_publication WHERE pubname = $1)", publicationName).Scan(&exists); err != nil {
		return errors.Raise(errors.ErrPostgresPublicationFailed, "failed to fetch publication", err)
	}

	if !exists {
		if _, err := s.sqlConn.Exec(
			ctx,
			fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE %s", publicationName, s.getTables()),
		); err != nil {
			return errors.Raise(errors.ErrPostgresPublicationFailed, "failed to create publication", err)
		}
	}

	return nil
}

func (s *pgSource) ensureReplicationSlot(ctx context.Context) error {
	_, err := pglogrepl.CreateReplicationSlot(
		ctx, s.repl, slotName, pluginName, pglogrepl.CreateReplicationSlotOptions{
			Temporary: false,
			Mode:      pglogrepl.LogicalReplication,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if goerrors.As(err, &pgErr) {
			if pgErr.Code != errors.PgDuplicateObjectCode {
				return errors.Raise(errors.ErrPostgresReplicationFailed, "failed to create replication slot", err)
			}
		} else {
			return errors.Raise(errors.ErrPostgresReplicationFailed, "failed to create replication slot", err)
		}
	}

	return nil
}

func (s *pgSource) sendStandby(ctx context.Context, span trace.Span, log *zap.Logger) {
	if err := pglogrepl.SendStandbyStatusUpdate(ctx, s.repl,
		pglogrepl.StandbyStatusUpdate{
			WALWritePosition: s.lastLSN,
			WALFlushPosition: s.lastLSN,
			WALApplyPosition: s.lastLSN,
		},
	); err != nil {
		err = errors.Raise(errors.ErrPostgresStandbyFailed, "", err)
		errors.HandleSpanError(ctx, span, log, "failed to send standby status update to postgres", err)
		return
	}

	log.Debug("sent standby status update", zap.String("lsn", s.lastLSN.String()))
}

func (s *pgSource) getTables() string {
	tables := make([]string, 0, len(tableHandlerMap))
	for key := range tableHandlerMap {
		tables = append(tables, fmt.Sprintf("\"%s\"", key))
	}

	return strings.Join(tables, ", ")
}

func extractColumns(rel *pglogrepl.RelationMessage, tuple *pglogrepl.TupleData) map[string]any {
	result := make(map[string]any)

	for i, col := range rel.Columns {
		val := tuple.Columns[i]

		switch val.DataType {
		case 'n': // NULL
			result[col.Name] = nil
		case 'u': // unchanged TOAST
			continue
		case 't': // text format
			result[col.Name] = string(val.Data)
		}
	}

	return result
}
