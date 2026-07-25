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

	swarmID := "swarm-rolecollapsetest-1"
	taskID := "task-1"

	// worker completes 8 tasks, planner and reviewer complete 1 each.
	// worker's share: 8/10 = 80%, well over our 70% threshold.
	agentCounts := map[string]int{
		"worker":   8,
		"planner":  1,
		"reviewer": 1,
	}

	for agentID, count := range agentCounts {
		for i := 0; i < count; i++ {
			content := fmt.Sprintf("completed subtask %d", i)
			recipient := "planner"
			event := schema.NewEvent(swarmID, agentID, schema.EventMessage, schema.Payload{
				Content:          &content,
				RecipientAgentID: &recipient,
			}, schema.Metadata{Framework: "custom", SDKVersion: "0.1.0"})
			event.ParentTaskID = &taskID
			if err := producer.Publish(ctx, event); err != nil {
				log.Fatalf("publish failed: %v", err)
			}
		}
	}

	log.Println("role collapse test events published")
}
