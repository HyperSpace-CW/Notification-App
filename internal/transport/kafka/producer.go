package kafka

import (
	"context"
	"encoding/json"
	"github.com/HyperSpace-CW/Notification-App/internal/services"
	"github.com/segmentio/kafka-go"
	"log"
)

func StartEmailConsumer(brokers []string, topic string, groupID string, service services.NotificationService) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var notif EmailNotification
		if err := json.Unmarshal(m.Value, &notif); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		if err := service.SendCodeToEmail(notif.Email, notif.Code); err != nil {
			log.Printf("Send email failed: %v", err)
		} else {
			log.Printf("Email sent to %s", notif.Email)
		}
	}
}
