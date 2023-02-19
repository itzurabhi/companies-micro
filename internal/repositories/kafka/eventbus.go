package kafka

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/itzurabhi/companies-micro/internal/repositories"
	"github.com/sirupsen/logrus"
)

type companiesEventBus struct {
	producer *kafka.Producer
	topic    string
}

// PostEvent implements repositories.EventBus
func (bus companiesEventBus) PostEvent(_ context.Context, items ...interface{}) error {
	for _, item := range items {
		jsonData, err := json.Marshal(item)

		if err != nil {
			logrus.Error("serialising event failed", item)
			continue
		}

		if err := bus.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &bus.topic, Partition: kafka.PartitionAny},
			Value:          jsonData,
		}, nil); err != nil {
			return err
		}
	}

	return nil
}

func CreateCompaniesEventBus(producer *kafka.Producer, topic string) repositories.EventBus {
	return companiesEventBus{
		producer,
		topic,
	}
}
