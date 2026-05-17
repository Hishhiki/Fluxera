package events

import (
	"context"
	"time"

	"fluxera/internal/models"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers []string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, event models.Event) error {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}

	topic, err := TopicForEvent(event)
	if err != nil {
		return err
	}

	value, err := SerializeEvent(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(event.ID),
		Value: value,
	})

}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}
