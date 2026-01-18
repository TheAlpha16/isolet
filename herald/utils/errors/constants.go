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

	// Cache errors
	ErrCacheCallFailed ErrorCode = "CACHE-01"

	// Kafka errors
	ErrKafkaProducerCreationFailed ErrorCode = "KAFKA-01"
)

var msgMap = map[ErrorCode]string{
	// System errors
	ErrInternalError: "internal server error",

	// Fact errors
	ErrFactInvalidKey:       "fact has invalid key",
	ErrFactInvalidType:      "fact has invalid type",
	ErrFactInvalidOccuredAt: "fact has invalid occured at timestamp",

	// Kafka errors
	ErrKafkaProducerCreationFailed: "failed to create kafka producer",
}
