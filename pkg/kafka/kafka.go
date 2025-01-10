package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type KafkaConsumer struct {
	brokers       []string
	topic         string
	group         string
	consumerGroup sarama.ConsumerGroup
}

type ConsumerHandler struct {
	ready chan bool
}

// NewKafkaConsumer creates a new KafkaConsumer.
func NewKafkaConsumer(brokers []string, topic, group string) (*KafkaConsumer, error) {
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
		brokers:       brokers,
		topic:         topic,
		group:         group,
		consumerGroup: consumerGroup,
	}, nil
}

// Start begins message consumption from Kafka.
func (kc *KafkaConsumer) Start(ctx context.Context, handler sarama.ConsumerGroupHandler) {
	go func() {
		for {
			if err := kc.consumerGroup.Consume(ctx, []string{kc.topic}, handler); err != nil {
				log.Printf("Error during consumption: %v\n", err)
			}

			// Check if context was cancelled
			if ctx.Err() != nil {
				return
			}
		}
	}()
}

// Close shuts down the Kafka connection.
func (kc *KafkaConsumer) Close() error {
	return kc.consumerGroup.Close()
}

// Setup is run before the consumer starts consuming messages.
func (ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup is called after message consumption has ended.
func (ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processes messages from Kafka.
func (h ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message received: topic=%s partition=%d offset=%d key=%s value=%s",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		session.MarkMessage(message, "")
	}
	return nil
}

// RunConsumer starts Kafka message consumption and handles graceful shutdown signals.
func RunConsumer(brokers []string, topic, group string) {
	consumer, err := NewKafkaConsumer(brokers, topic, group)
	if err != nil {
		log.Fatalf("Failed to create KafkaConsumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	handler := ConsumerHandler{
		ready: make(chan bool),
	}

	consumer.Start(ctx, handler)

	// Handle system signals for graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	log.Println("KafkaConsumer started and running...")
	select {
	case <-sigchan:
		log.Println("Shutdown signal received")
		cancel()
	}

	if err := consumer.Close(); err != nil {
		log.Printf("Error closing KafkaConsumer: %v\n", err)
	}
}
