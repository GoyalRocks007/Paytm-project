package kafkaclient

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

var (
	GetKafkaClient = func(brokerURL, topic, groupID string) IKafkaClient {
		return &KafkaClient{
			reader: kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{brokerURL},
				Topic:   topic,
				GroupID: groupID,
			}),
		}
	}
)

type KafkaClient struct {
	reader *kafka.Reader
}

func (c *KafkaClient) Listen(ctx context.Context, handler func(key, value []byte)) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println("❌ Error fetching message:", err)
			continue
		}
		log.Printf("📨 Received message: key=%s value=%s", string(msg.Key), string(msg.Value))
		handler(msg.Key, msg.Value)

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Println("❌ Commit error:", err)
		}
	}
}

func (c *KafkaClient) Stop() error {
	return c.reader.Close()
}
