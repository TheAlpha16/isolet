package consumer

import (
	"context"
	"time"

	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	factDom "github.com/TheAlpha16/isolet/oracle/internal/domain/fact"
	"github.com/TheAlpha16/isolet/oracle/utils"
	"github.com/TheAlpha16/isolet/oracle/utils/tracer"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var consumerTracer = otel.Tracer("oracle.consumer")

type Deserializer func(context.Context, []byte) (factDom.Fact, error)

type Consumer struct {
	consumer      *kafka.Consumer
	factUc        factDom.Usecase
	deserializers map[string]Deserializer
}

func (c *Consumer) Start(ctx context.Context) error {
	ctx = errorDom.SetSrcInCtx(ctx, "oracle.consumer")
	var topics []string
	for topic := range c.deserializers {
		topics = append(topics, topic)
	}

	err := c.consumer.SubscribeTopics(topics, nil)
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrKafkaSubscriptionFailed, "failed to subscribe to topics", err, nil)
	}

	for {
		select {
		case <-ctx.Done():
			return c.consumer.Close()
		case ev := <-c.consumer.Events():
			eventCtx, eventSpan, logger := tracer.StartSpan(ctx, consumerTracer, "oracle.consumer.handler")
			defer eventSpan.End()

			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				deserializer, ok := c.deserializers[*e.TopicPartition.Topic]
				if !ok {
					tracer.FactsTotal.WithLabelValues(*e.TopicPartition.Topic, "ignored", "unknown").Inc()
					logger.Warn("no deserializer found for topic", zap.String("topic", *e.TopicPartition.Topic))
					continue
				}

				start := time.Now()
				fact, err := deserializer(eventCtx, e.Value)
				if err != nil {
					tracer.FactsTotal.WithLabelValues(*e.TopicPartition.Topic, tracer.StatusFailure, "unknown").Inc()
					errorDom.RaiseToSentry(ctx, err)
					logger.Error("failed to deserialize fact", zap.Error(err))
					continue
				}

				err = c.factUc.HandleEvent(eventCtx, fact)
				status := tracer.StatusSuccess
				if err != nil {
					status = tracer.StatusFailure
					logger.Error("failed to handle event", zap.Error(err), zap.String("fact_type", string(fact.FactType())))
				}

				tracer.FactsTotal.WithLabelValues(*e.TopicPartition.Topic, status, string(fact.FactType())).Inc()
				tracer.FactProcessingDurationSeconds.WithLabelValues(*e.TopicPartition.Topic, string(fact.FactType())).Observe(time.Since(start).Seconds())
				tracer.FactDeliveryDurationSeconds.WithLabelValues(*e.TopicPartition.Topic, string(fact.FactType())).Observe(start.Sub(fact.OccuredAt()).Seconds())

			case kafka.Error:
				if e.IsFatal() {
					err := errorDom.Raise(eventCtx, errorDom.ErrKafkaConsumerError, "fatal kafka error", e, nil)
					logger.Error("fatal kafka error", zap.Error(err))
					return err
				}
				logger.Error("kafka error", zap.Error(e))
			default:
				logger.Debug("ignored kafka event", zap.Any("event", e))
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

func New(ctx context.Context, factUc factDom.Usecase) (*Consumer, error) {
	config := utils.GetConfig()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        config.Kafka.Brokers,
		"group.id":                 config.Kafka.GroupID,
		"auto.offset.reset":        "earliest",
		"go.events.channel.enable": true,
	})

	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrKafkaConnectionFailed, "failed to create kafka consumer", err, nil)
	}

	return &Consumer{
		consumer: c,
		factUc:   factUc,
		deserializers: map[string]Deserializer{
			config.Instances.FactTopic: deserializeInstanceFact,
		},
	}, nil
}
