package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/TheAlpha16/isolet/api/utils"
	"github.com/TheAlpha16/isolet/api/utils/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type tablesMap map[string]TableInfo

const STRUCT_DB_TAG = "db"

var Tables = make(tablesMap) // Tables

type BaseModel struct {
	ID        int64 `db:"id"`
	CreatedAt int64 `db:"created_at"`
	UpdatedAt int64 `db:"updated_at"`
}

func (t tablesMap) Register(tableName string, tab interface{}) {
	if _, ok := t[tableName]; ok {
		return
	}
	fieldNames := utils.StructTagsAsString(tab, STRUCT_DB_TAG, 1)
	var fieldNamesWithPrefix string
	fieldNameSlice := strings.SplitSeq(fieldNames, ",")
	for f := range fieldNameSlice {
		fieldNamesWithPrefix += fmt.Sprintf("%s.%s,", tableName, f)
	}
	if len(fieldNamesWithPrefix) > 0 {
		fieldNamesWithPrefix = fieldNamesWithPrefix[:len(fieldNamesWithPrefix)-1]
	}
	t[tableName] = TableInfo{
		FieldNames:           fieldNames,
		FieldNamesWithPrefix: fieldNamesWithPrefix,
	}
}

type TableInfo struct {
	FieldNames           string
	FieldNamesWithPrefix string
}

type myQueryTracer struct {
	config *utils.Config
	log    *logger.StandardLogger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	ctx, _ = otel.Tracer("pgx").Start(ctx, "PGXQuery", trace.WithAttributes(
		attribute.String("db.system", "Postgres"),
		attribute.String("db.statement", data.SQL),
		attribute.String("db.name", conn.Config().Database),
	))
	logger.GetLogger(ctx).Debug(
		"Executing SQL command",
		zap.String("db", conn.Config().Database),
		zap.String("sql", data.SQL),
		zap.Any("args", data.Args),
	)
	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	if data.Err != nil {
		span.RecordError(data.Err)
	}
	span.End()
}

func NewConnection(ctx context.Context, dbURI string) (*pgxpool.Pool, func(), error) {
	config := utils.GetConfig()
	logger := logger.GetAppLogger()
	pgxConfig, err := pgxpool.ParseConfig(dbURI)
	if err != nil {
		return nil, nil, err
	}
	pgxConfig.MaxConns = int32(config.Database.MaxConnections)
	pgxConfig.MaxConnLifetime = config.Database.MaxConnectionIdleTime
	pgxConfig.ConnConfig.Tracer = &myQueryTracer{log: logger, config: config}
	dbPool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, nil, err
	}
	closeConn := func() {
		dbPool.Close()
	}
	return dbPool, closeConn, nil
}
