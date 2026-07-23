package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shylin26/swarmlens/internal/ingest"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	consumer := ingest.NewConsumer([]string{"localhost:19092"}, "agent.messages", "swarmlens-control-plane")
	defer consumer.Close()

	log.Println("swarmlens control plane starting, listening on agent.messages...")

	for {
		event, err := consumer.ReadEvent(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("shutdown signal received, exiting cleanly")
				return
			}
			log.Printf("error reading event: %v", err)
			continue
		}

		log.Printf("received event: id=%s swarm=%s agent=%s type=%s",
			event.EventID, event.SwarmID, event.AgentID, event.Type)
	}
}
