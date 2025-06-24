package consumer

import (
	"context"
	"github.com/HyperSpace-CW/Notification-App/internal/services"
	pb "github.com/HyperSpace-CW/Notification-App/pkg/proto"
	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"
	"log"
)

type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	topic         string
	service       services.NotificationService
}

func NewKafkaConsumer(brokers []string, topic, group string, service services.NotificationService) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Version = sarama.V2_8_0_0

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		topic:         topic,
		service:       service,
	}, nil
}

func (kc *KafkaConsumer) Start(ctx context.Context) {
	handler := &ConsumerHandler{
		service: kc.service,
	}

	go func() {
		for {
			if err := kc.consumerGroup.Consume(ctx, []string{kc.topic}, handler); err != nil {
				log.Printf("Error during consumption: %v", err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumerGroup.Close()
}

// ConsumerHandler реализует интерфейс sarama.ConsumerGroupHandler
type ConsumerHandler struct {
	service services.NotificationService
}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// Вызывается перед началом обработки сообщений
	return nil
}

func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	// Вызывается после завершения обработки сообщений
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg pb.ConfirmRegistrationRequest
		if err := proto.Unmarshal(message.Value, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// Вызов NotificationService
		if err := h.service.SendCodeToEmail(context.Background(), msg.Email, msg.Code); err != nil {
			log.Printf("Failed to send email: %v", err)
		}

		// Сообщить Kafka, что сообщение обработано
		session.MarkMessage(message, "")
	}
	return nil
}
