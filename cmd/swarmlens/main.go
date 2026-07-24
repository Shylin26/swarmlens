package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shylin26/swarmlens/internal/ingest"
	"github.com/Shylin26/swarmlens/internal/metrics"
	"github.com/Shylin26/swarmlens/internal/state"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	metrics.StartServer(":2112")

	consumer := ingest.NewConsumer([]string{"localhost:19092"}, "agent.messages", "swarmlens-control-plane")
	defer consumer.Close()

	store := state.NewStore("localhost:16379")

	log.Println("swarmlens control plane starting, listening on agent.messages...")
	log.Println("metrics available at http://localhost:2112/metrics")

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

		metrics.EventsTotal.WithLabelValues(event.SwarmID, event.AgentID, string(event.Type)).Inc()

		if err := store.RecordAgentActivity(ctx, event.SwarmID, event.AgentID); err != nil {
			log.Printf("error recording agent activity: %v", err)
			continue
		}
	}
}
