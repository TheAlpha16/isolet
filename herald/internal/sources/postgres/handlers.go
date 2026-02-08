package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/herald/pkg/facts"
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/tracer"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgproto3"
)

func (s *pgSource) handleReplication(ctx context.Context, out chan<- facts.Fact) error {
	ctx, span, log := tracer.StartSpan(ctx, pgTracer, "herald.sources.postgres.handleReplication")
	defer span.End()

	standbyTimeout := time.NewTicker(utils.GetConfig().Postgres.StandbyTimeout)
	defer standbyTimeout.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-standbyTimeout.C:
			// periodic standby update (heartbeat)
			if s.lastLSN > 0 {
				s.sendStandby(ctx, span, log)
			}

		default:
			msg, err := s.repl.ReceiveMessage(ctx)
			if err != nil {
				return err
			}

			switch m := msg.(type) {
			case *pgproto3.CopyData:
				// handle replication messages in the CopyData response

				switch m.Data[0] {
				case pglogrepl.XLogDataByteID:
					logData, err := pglogrepl.ParseXLogData(m.Data[1:])
					if err != nil {
						continue
					}
					s.lastLSN = logData.WALStart + pglogrepl.LSN(len(logData.WALData))

					fmt.Printf("Received WAL data at LSN %s\n", logData.WALStart)
					fmt.Printf("Raw: %v\n", logData.WALData)

					s.sendStandby(ctx, span, log)

				case pglogrepl.PrimaryKeepaliveMessageByteID:
					keepAliveMsg, err := pglogrepl.ParsePrimaryKeepaliveMessage(m.Data[1:])
					if err != nil {
						continue
					}
					if keepAliveMsg.ReplyRequested {
						s.sendStandby(ctx, span, log)
						standbyTimeout.Reset(utils.GetConfig().Postgres.StandbyTimeout)
					}
				}
			default:
			}
		}
	}
}
