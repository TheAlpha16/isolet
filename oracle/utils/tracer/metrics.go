package tracer

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	StatusSuccess = "success"
	StatusFailure = "failure"
)

var (
	// Challenges
	FlagSubmissionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flag_submissions_total",
			Help: "Total number of flag submissions",
		},
		[]string{"challenge_id", "correct"},
	)

	// Instances
	InstanceOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "instance_operations_total",
			Help: "Total number of instance operations",
		},
		[]string{"challenge_id", "operation", "status"},
	)

	InstanceActiveTotal = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "instance_active_total",
			Help: "Current number of active instances",
		},
		[]string{"challenge_id"},
	)

	InstanceProvisionDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "instance_provision_duration_seconds",
			Help: "Time taken to provision an instance",
		},
		[]string{"challenge_id", "operation", "status"},
	)

	// Emails
	EmailsSentTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "emails_sent_total",
			Help: "Total number of emails sent",
		},
		[]string{"status"},
	)

	EmailSendDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "email_send_duration_seconds",
			Help: "Time taken to send email",
		},
		[]string{"status"},
	)

	// HTTP
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	HttpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests",
		},
		[]string{"method", "route"},
	)

	HttpRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of in-flight HTTP requests",
		},
		[]string{"route"},
	)

	// Facts
	FactsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "facts_total",
			Help: "Total number of facts processed",
		},
		[]string{"topic", "status", "fact_type"},
	)

	FactDeliveryDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "fact_delivery_duration_seconds",
			Help: "Time taken to deliver a fact",
		},
		[]string{"topic", "fact_type"},
	)

	FactProcessingDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "fact_processing_duration_seconds",
			Help: "Time taken to process a fact",
		},
		[]string{"topic", "fact_type"},
	)
)

func InitMetrics() error {
	return nil
}
