package producer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer interface {
	SendMovieEvent(movieBytes []byte) (err error)
}

type KafkaProducerConfig struct {
	Producer *kafka.Producer
	Topic    *string
}

func (k *KafkaProducerConfig) SendMovieEvent(movieBytes []byte) error {
	err := k.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: k.Topic, Partition: kafka.PartitionAny},
		Value:          movieBytes,
	}, nil)

	if err != nil {
		return fmt.Errorf("error sending event to kafka %v", err)
	}

	return nil
}

func GetKafkaProducer(topic *string) (producer KafkaProducerConfig, err error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	if err != nil {
		panic(err)
	}

	producer = KafkaProducerConfig{Producer: p, Topic: topic}

	return
}
