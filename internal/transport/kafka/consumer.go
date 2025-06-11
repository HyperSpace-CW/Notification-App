package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

func SendEmailNotification(brokers []string, topic string, message EmailNotification) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	defer writer.Close()

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: data,
		})
}
