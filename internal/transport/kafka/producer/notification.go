package producer

import (
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (kp *KafkaProducer) SendMessage(message proto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(data),
	}

	_, _, err = kp.producer.SendMessage(msg)
	return err
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
