package kafka

import (
	"context"

	"github.com/TheAlpha16/isolet/herald/internal/emitter"
	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/TheAlpha16/isolet/herald/utils/kafka"
	"github.com/TheAlpha16/isolet/herald/utils/logger"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var tracer = otel.Tracer("herald.emit.kafka")

type kafkaEmitter struct {
	topic    string
	producer *kafka.Client
}

func (e *kafkaEmitter) Emit(ctx context.Context, fact facts.Fact) error {
	ctx, span := tracer.Start(ctx, "herald.emit.kafka")
	defer span.End()

	log := logger.GetAppLogger().With(
		zap.String("emitter", "kafka"),
		zap.ByteString("key", fact.Key()),
		zap.String("fact_type", string(fact.FactType())),
	)
	log.Debug("received fact for emission")

	value, err := fact.Marshal()
	if err != nil {
		errors.HandleSpanError(ctx, span, log, "failed to marshal fact", err)
		return err
	}

	log.Debug(
		"producing fact to kafka",
		zap.ByteString("payload", value),
	)

	if err := e.producer.Produce(e.topic, fact.Key(), value); err != nil {
		err = errors.Raise(errors.ErrKafkaProduceFailed, "", err)
		errors.HandleSpanError(ctx, span, log, "failed to produce fact to kafka", err)
		return err
	}

	log.Debug("produced fact to kafka successfully")
	return nil
}

func (e *kafkaEmitter) Close() error {
	e.producer.Close()
	return nil
}

func NewEmitter(topic string, producer *kafka.Client) emitter.Emitter {
	return &kafkaEmitter{
		topic:    topic,
		producer: producer,
	}
}
