package events

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventProducer interface {
	Publish(ctx context.Context, topic string, key string, value []byte) error
	Close() error
}

type Producer struct {
	Writer *kafka.Writer
}

var _ EventProducer = (*Producer)(nil)

func NewProducer(brokers []string) *Producer {
	return &Producer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    1024,
			Async:        true,
			BatchTimeout: 10 * time.Second,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, value []byte) error {
	err := p.Writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
	if err != nil {
		log.Printf("❌ Failed to publish message to topic %s: %v", topic, err)
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	if err := p.Writer.Close(); err != nil {
		log.Printf("❌ Failed to close Kafka producer: %v", err)
		return err
	}
	return nil
}
