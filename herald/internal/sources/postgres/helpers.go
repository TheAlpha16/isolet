package postgres

import (
	"context"
	goerrors "errors"
	"fmt"

	"github.com/TheAlpha16/isolet/herald/utils/errors"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *pgSource) ensurePublication(ctx context.Context) error {
	var exists bool

	if err := s.sqlConn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_publication WHERE pubname = $1)", publicationName).Scan(&exists); err != nil {
		return errors.Raise(errors.ErrPostgresPublicationFailed, "failed to fetch publication", err)
	}

	if !exists {
		if _, err := s.sqlConn.Exec(ctx, fmt.Sprintf("CREATE PUBLICATION %s", publicationName)); err != nil {
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
