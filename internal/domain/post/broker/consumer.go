package broker

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/izzanzahrial/skeleton/config"
)

func NewConsumer() (sarama.ConsumerGroup, error) {
	cfg, err := config.NewConsumer()
	if err != nil {
		return nil, fmt.Errorf("failed to get kafka consumer config: %w", err)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V3_3_0_0
	config.Consumer.Offsets.AutoCommit.Enable = cfg.AutoCommit      // default true, use false higher level of control, ensuring that messages are processed at least once or exactly once, depending on the use case
	config.Consumer.Fetch.Default = cfg.FetchBytes * cfg.FetchBytes // 1MB, increase batch size to minimize round-trip when fetching messages
	config.Consumer.MaxWaitTime = cfg.MaxWait                       // reduce the numbers of round-trip so the consumer doesn't have to keep asking when there's no data
	config.ChannelBufferSize = cfg.ChannelBufferSize                // default 256, increase buffer size if consumer is slow, but at the cost of memory usage
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(cfg.Addresses, cfg.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return consumer, nil
}

type Handler struct{}

func (consumer Handler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (consumer Handler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (consumer Handler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// process the message
		fmt.Printf("Message topic:%q partition:%d offset:%d value:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		// after processing the message, mark the offset
		sess.MarkMessage(msg, "")
	}
	return nil
}
