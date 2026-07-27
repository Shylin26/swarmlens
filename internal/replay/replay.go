package replay

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/Shylin26/swarmlens/internal/schema"
)

type Replayer struct {
	brokers []string
	topic   string
}

func NewReplayer(brokers []string, topic string) *Replayer {
	return &Replayer{brokers: brokers, topic: topic}

}

func (r *Replayer) Replay(ctx context.Context, swarmID string, fromTime time.Time, partitionCount int, handler func(schema.Event)) error {
	for partition := 0; partition < partitionCount; partition++ {
		if err := r.replayPartition(ctx, swarmID, fromTime, partition, handler); err != nil {
			return fmt.Errorf("replay: partition %d: %w", partition, err)

		}
	}
	return nil
}

func (r *Replayer) replayPartition(ctx context.Context, swarmID string, fromTime time.Time, partition int, handler func(schema.Event)) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   r.brokers,
		Topic:     r.topic,
		Partition: partition,
	})
	defer reader.Close()

	if err := reader.SetOffsetAt(ctx, fromTime); err != nil {
		return fmt.Errorf("failed to seek to time: %w", err)
	}
	lastOffset, err := reader.ReadLag(ctx)
	if err != nil {
		return fmt.Errorf("failed to read lag: %w", err)

	}
	if lastOffset == 0 {
		return nil
	}
	for i := int64(0); i < lastOffset; i++ {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var event schema.Event
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			continue
		}

		if event.SwarmID == swarmID {
			handler(event)
		}
	}

	return nil

}
func (r *Replayer) CollectEvents(ctx context.Context, swarmID string, from time.Time, partitionCount int) ([]schema.Event, error) {
	var events []schema.Event
	err := r.Replay(ctx, swarmID, from, partitionCount, func(event schema.Event) {
		events = append(events, event)
	})
	if err != nil {
		return nil, err
	}
	return events, nil
}
