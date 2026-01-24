package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var (
	// Challenges
	FlagSubmissionsTotal metric.Int64Counter

	// Instances
	InstanceOperationsTotal          metric.Int64Counter
	InstanceActiveTotal              metric.Int64UpDownCounter
	InstanceProvisionDurationSeconds metric.Float64Histogram

	// Emails
	EmailsSentTotal          metric.Int64Counter
	EmailSendDurationSeconds metric.Float64Histogram
)

func InitMetrics(ctx context.Context) error {
	exporter, err := prometheus.New(
		prometheus.WithoutScopeInfo(),
	)
	if err != nil {
		return err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	// Initialize instruments
	meter := otel.Meter("oracle")

	FlagSubmissionsTotal, err = meter.Int64Counter(
		"flag_submissions_total",
		metric.WithDescription("Total number of flag submissions"),
	)
	if err != nil {
		return err
	}

	InstanceOperationsTotal, err = meter.Int64Counter(
		"instance_operations_total",
		metric.WithDescription("Total number of instance operations"),
	)
	if err != nil {
		return err
	}

	InstanceActiveTotal, err = meter.Int64UpDownCounter(
		"instance_active_total",
		metric.WithDescription("Current number of active instances"),
	)
	if err != nil {
		return err
	}

	InstanceProvisionDurationSeconds, err = meter.Float64Histogram(
		"instance_provision_duration_seconds",
		metric.WithDescription("Time taken to provision an instance"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(1, 2, 5, 10, 15, 20, 30, 60),
	)
	if err != nil {
		return err
	}

	EmailsSentTotal, err = meter.Int64Counter(
		"emails_sent_total",
		metric.WithDescription("Total number of emails sent"),
	)
	if err != nil {
		return err
	}

	EmailSendDurationSeconds, err = meter.Float64Histogram(
		"email_send_duration_seconds",
		metric.WithDescription("Time taken to send email"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	return nil
}
