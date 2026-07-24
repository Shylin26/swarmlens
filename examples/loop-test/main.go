package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Shylin26/swarmlens/internal/ingest"
	"github.com/Shylin26/swarmlens/internal/schema"
)

func main() {
	ctx := context.Background()
	producer := ingest.NewProducer([]string{"localhost:19092"}, "agent.messages")

	swarmID := "swarm-looptest-1"
	taskID := "task-1"

	for i := 0; i < 4; i++ {
		content := fmt.Sprintf("ping %d", i)
		recipient := "reviewer"
		event := schema.NewEvent(swarmID, "worker", schema.EventMessage, schema.Payload{
			Content:          &content,
			RecipientAgentID: &recipient,
		}, schema.Metadata{Framework: "custom", SDKVersion: "0.1.0"})
		event.ParentTaskID = &taskID
		if err := producer.Publish(ctx, event); err != nil {
			log.Fatalf("publish failed: %v", err)
		}

		content2 := fmt.Sprintf("pong %d", i)
		recipient2 := "worker"
		event2 := schema.NewEvent(swarmID, "reviewer", schema.EventMessage, schema.Payload{
			Content:          &content2,
			RecipientAgentID: &recipient2,
		}, schema.Metadata{Framework: "custom", SDKVersion: "0.1.0"})
		event2.ParentTaskID = &taskID
		if err := producer.Publish(ctx, event2); err != nil {
			log.Fatalf("publish failed: %v", err)
		}
	}

	log.Println("loop test events published")
}
