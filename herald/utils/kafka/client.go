package kafka

import (
	"github.com/TheAlpha16/isolet/herald/utils"
	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Client struct {
	producer *kafka.Producer
}

func (c *Client) Produce(topic string, key, value []byte) error {
	return c.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: value,
	}, nil)
}

func (c *Client) Close() {
	c.producer.Flush(5000)
	c.producer.Close()
}

func NewClient() (*Client, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": utils.GetConfig().Kafka.Brokers,
		"retries":           utils.GetConfig().Kafka.Retries,
	})
	if err != nil {
		return nil, errors.Raise(errors.ErrKafkaProducerCreationFailed, "", err)
	}

	return &Client{producer: p}, nil
}
