package broker

import (
	"context"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(producer sarama.SyncProducer) *Producer {
	return &Producer{producer: producer}
}

func (p *Producer) Publish(ctx context.Context, topic string, payload []byte) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(payload),
	}

	_, _, err := p.producer.SendMessage(message)
	return err
}
