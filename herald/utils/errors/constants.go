package errors

const (
	SVC string = "HRLD"
)

const (
	// System errors
	ErrInternalError ErrorCode = "SYS-00"

	// Fact errors
	ErrFactInvalidKey       ErrorCode = "FACT-01"
	ErrFactInvalidType      ErrorCode = "FACT-02"
	ErrFactInvalidOccuredAt ErrorCode = "FACT-03"
	ErrFactInvalidField     ErrorCode = "FACT-04"
	ErrFactMarshalFailed    ErrorCode = "FACT-05"

	// Cache errors
	ErrCacheCallFailed ErrorCode = "CACHE-01"

	// Kafka errors
	ErrKafkaProducerCreationFailed ErrorCode = "KAFKA-01"
	ErrKafkaProduceFailed          ErrorCode = "KAFKA-02"

	// K8s errors
	ErrK8sInformerCreationFailed ErrorCode = "K8S-01"
	ErrK8sCacheFailed            ErrorCode = "K8S-02"

	// Postgres errors
	ErrPostgresConnectionFailed  ErrorCode = "PG-01"
	ErrPostgresPublicationFailed ErrorCode = "PG-02"
	ErrPostgresReplicationFailed ErrorCode = "PG-03"
	ErrPostgresStandbyFailed     ErrorCode = "PG-04"
	ErrPostgresRelationMissing   ErrorCode = "PG-05"
	ErrPostgresHandlerMissing    ErrorCode = "PG-06"
	ErrPostgresInvalidValue      ErrorCode = "PG-07"
)

var msgMap = map[ErrorCode]string{
	// System errors
	ErrInternalError: "internal server error",

	// Fact errors
	ErrFactInvalidKey:       "fact has invalid key",
	ErrFactInvalidType:      "fact has invalid type",
	ErrFactInvalidOccuredAt: "fact has invalid occurred at timestamp",

	// Kafka errors
	ErrKafkaProducerCreationFailed: "failed to create kafka producer",
	ErrKafkaProduceFailed:          "failed to produce kafka message",

	// K8s errors
	ErrK8sInformerCreationFailed: "failed to create k8s informer",

	// Postgres errors
	ErrPostgresConnectionFailed: "failed to connect to postgres",
	ErrPostgresStandbyFailed:    "failed to send standby status update to postgres",
	ErrPostgresRelationMissing:  "unknown relation in replication message",
	ErrPostgresHandlerMissing:   "no handler for table in replication message",
}

const (
	PgDuplicateObjectCode = "42710"
)
