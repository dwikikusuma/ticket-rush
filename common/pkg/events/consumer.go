package events

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventConsumer interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type Consumer struct {
	Reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID string, topic string) *Consumer {
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			GroupID:     groupID,
			Topic:       topic,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			MaxWait:     1 * time.Second,
			StartOffset: kafka.FirstOffset,
		}),
	}
}

func (c *Consumer) Close() error {
	return c.Reader.Close()
}

func (c *Consumer) CommitMessages(ctx context.Context, msg ...kafka.Message) error {
	return c.Reader.CommitMessages(ctx, msg...)
}

func (c *Consumer) FetchMessage(ctx context.Context) (kafka.Message, error) {
	return c.Reader.FetchMessage(ctx)
}
