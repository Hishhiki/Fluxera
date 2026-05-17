package events

import (
	"context"
	"encoding/json"
	"fluxera/internal/models"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler Handler
}

func NewKafkaConsumer(brokers []string, topics []string, groupID string, handler Handler) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupTopics: topics,
			GroupID:     groupID,
		}),
		handler: handler,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	go func() {
		for {
			message, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}

				log.Printf("failed to fetch kafka message: %v", err)
				time.Sleep(time.Second)
				continue
			}

			var event models.Event
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("failed to unmarshal kafka event: %v", err)
				continue
			}

			if err := c.handler.Handle(ctx, event); err != nil {
				log.Printf("failed to handle kafka event: %v", err)
				continue
			}

			if err := c.reader.CommitMessages(ctx, message); err != nil {
				log.Printf("failed to commit kafka message: %v", err)
				continue
			}
		}
	}()
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
