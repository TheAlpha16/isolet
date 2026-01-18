package kafka

import (
	"context"

	"github.com/TheAlpha16/isolet/herald/internal/emitter"
	"github.com/TheAlpha16/isolet/herald/internal/facts"
	"github.com/TheAlpha16/isolet/herald/utils/kafka"
)

type kafkaEmitter struct {
	topic    string
	producer *kafka.Client
}

func (e *kafkaEmitter) Emit(_ context.Context, fact facts.Fact) error {
	value, err := fact.Marshal()
	if err != nil {
		return err
	}

	return e.producer.Produce(e.topic, fact.Key(), value)
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
