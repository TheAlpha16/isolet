package tracer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	StatusSuccess = "success"
	StatusFailure = "failure"
	StatusUnknown = "unknown"

	OpStart  = "start"
	OpStop   = "stop"
	OpExtend = "extend"
	OpExpire = "expire"
)

var (
	LabelChallengeID = "challenge_id"
	LabelCorrect     = "correct"
	LabelOperation   = "operation"
	LabelStatus      = "status"
	LabelMethod      = "method"
	LabelRoute       = "route"
	LabelTopic       = "topic"
	LabelFactType    = "fact_type"
)

var (
	// Challenges
	FlagSubmissionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flag_submissions_total",
			Help: "Total number of flag submissions",
		},
		[]string{LabelChallengeID, LabelCorrect},
	)

	// Instances
	InstanceOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instance_operations_total",
			Help: "Total number of instance operations",
		},
		[]string{LabelChallengeID, LabelOperation, LabelStatus},
	)

	InstanceProvisionDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "instance_provision_duration_seconds",
			Help: "Time taken to provision an instance",
		},
		[]string{LabelChallengeID, LabelOperation, LabelStatus},
	)

	// Emails
	EmailsSentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "emails_sent_total",
			Help: "Total number of emails sent",
		},
		[]string{LabelStatus},
	)

	EmailSendDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "email_send_duration_seconds",
			Help: "Time taken to send email",
		},
		[]string{LabelStatus},
	)

	// HTTP
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{LabelMethod, LabelRoute, LabelStatus},
	)

	HttpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{LabelMethod, LabelRoute},
	)

	HttpRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of in-flight HTTP requests",
		},
		[]string{LabelRoute},
	)

	// Facts
	FactsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "facts_total",
			Help: "Total number of facts processed",
		},
		[]string{LabelTopic, LabelStatus, LabelFactType},
	)

	FactDeliveryDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "fact_delivery_duration_seconds",
			Help: "Time taken to deliver a fact",
		},
		[]string{LabelTopic, LabelFactType},
	)

	FactProcessingDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "fact_processing_duration_seconds",
			Help: "Time taken to process a fact",
		},
		[]string{LabelTopic, LabelFactType},
	)
)

func InitMetrics() error {
	return nil
}
