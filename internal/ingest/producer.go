package ingest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shylin26/swarmlens/internal/schema"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
	return &Producer{writer: writer}
}
func (p *Producer) Publish(ctx context.Context, event schema.Event) error {
	if err := event.Validate(); err != nil {
		return fmt.Errorf("ingest:invalid event :%w", err)
	}
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("ingest: failed to marshal event: %w", err)
	}
	msg := kafka.Message{
		Key:   []byte(event.SwarmID),
		Value: value,
	}
	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("ingest:failed to write message :%w", err)
	}
	return nil
}
