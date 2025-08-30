package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writers map[string]*kafka.Writer
}

type MessagePayload struct {
	Type      string `json:"type"`
	Data      any    `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

func NewProducer(brokers []string, topics map[string]string) *Producer {
	writers := make(map[string]*kafka.Writer)
	for name, topic := range topics {
		writers[name] = &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}
	}

	return &Producer{
		writers: writers,
	}
}

func (p *Producer) SendMessage(ctx context.Context, writerName string, key string, payload MessagePayload) error {
	writer, exists := p.writers[writerName]
	if !exists {
		slog.Error("writer not found", "writerName", writerName)
		return nil
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		slog.Error("failed to marshal payload", "error", err)
		return nil
	}

	return writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: jsonPayload,
	})
}

func (p *Producer) Close() error {
	for _, writer := range p.writers {
		if err := writer.Close(); err != nil {
			slog.Error("Error closing writer", "error", err)
		}
	}
	return nil
}
