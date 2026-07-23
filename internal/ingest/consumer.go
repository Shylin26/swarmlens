package ingest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shylin26/swarmlens/internal/schema"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
	})
	return &Consumer{reader: reader}

}
func (c *Consumer) ReadEvent(ctx context.Context) (schema.Event, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return schema.Event{}, fmt.Errorf("ingest: failed to read message: %w", err)
	}
	var event schema.Event
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return schema.Event{}, fmt.Errorf("ingest: failed to unmarshal event: %w", err)
	}
	return event, nil

}
func (c *Consumer) Close() error {
	return c.reader.Close()
}
