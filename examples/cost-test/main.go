package main

import (
	"context"
	"log"

	"github.com/Shylin26/swarmlens/internal/ingest"
	"github.com/Shylin26/swarmlens/internal/schema"
)

func main() {
	ctx := context.Background()
	producer := ingest.NewProducer([]string{"localhost:19092"}, "agent.messages")

	swarmID := "swarm-costtest-1"

	// A few normal-cost events to establish a baseline.
	for i := 0; i < 3; i++ {
		content := "normal response"
		recipient := "reviewer"
		cost := 0.0005
		event := schema.NewEvent(swarmID, "worker", schema.EventMessage, schema.Payload{
			Content:          &content,
			RecipientAgentID: &recipient,
			CostUSD:          &cost,
		}, schema.Metadata{Framework: "custom", SDKVersion: "0.1.0"})
		if err := producer.Publish(ctx, event); err != nil {
			log.Fatalf("publish failed: %v", err)
		}
	}

	// One deliberately huge cost, well over 5x the ~0.0005 baseline.
	content := "expensive response"
	recipient := "reviewer"
	spikeCost := 0.01
	spikeEvent := schema.NewEvent(swarmID, "worker", schema.EventMessage, schema.Payload{
		Content:          &content,
		RecipientAgentID: &recipient,
		CostUSD:          &spikeCost,
	}, schema.Metadata{Framework: "custom", SDKVersion: "0.1.0"})
	if err := producer.Publish(ctx, spikeEvent); err != nil {
		log.Fatalf("publish failed: %v", err)
	}

	log.Println("cost test events published")
}
