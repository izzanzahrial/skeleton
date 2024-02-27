package broker

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/izzanzahrial/skeleton/config"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer() (*Producer, error) {
	cfg, err := config.NewProducer()
	if err != nil {
		return nil, fmt.Errorf("failed to get kafka producer config: %w", err)
	}

	kafka := sarama.NewConfig()
	kafka.Version = sarama.V3_3_0_0        // TODO: dynamic versioning
	kafka.Producer.Return.Successes = true // This is mandatory
	kafka.Producer.Timeout = cfg.Timeout
	kafka.ChannelBufferSize = cfg.ChannelBufferSize // default 256, increase buffer size if consumer is slow, but at the cost of memory usage

	// use this in production
	// kafkaCfg.Producer.Flush.Bytes = cfg.FlushBytes * cfg.FlushBytes // 1MB, increase batch size to minimize round-trip when sending messages

	kafkaProducer, err := sarama.NewSyncProducer(cfg.Addresses, kafka)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: kafkaProducer}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, payload []byte) error {
	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(payload),
	}

	_, _, err := p.producer.SendMessage(message)
	if err != nil {
		return err
	}

	return nil
}
