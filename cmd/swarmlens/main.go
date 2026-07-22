package main

import (
	"context"
	"log"

	"github.com/Shylin26/swarmlens/internal/ingest"
	"github.com/Shylin26/swarmlens/internal/schema"
)

func main() {
	producer := ingest.NewProducer([]string{"localhost:19092"}, "agent.messages")

	meta := schema.Metadata{Framework: "custom", SDKVersion: "0.1.0"}
	content := "hello from swarmlens"
	event := schema.NewEvent("swarm-test-1", "agent-1", schema.EventMessage, schema.Payload{
		Content: &content,
	}, meta)

	ctx := context.Background()
	if err := producer.Publish(ctx, event); err != nil {
		log.Fatalf("failed to publish event: %v", err)
	}

	log.Println("event published successfully")
}
