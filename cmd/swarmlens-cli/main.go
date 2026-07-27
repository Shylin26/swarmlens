package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shylin26/swarmlens/internal/replay"
	"github.com/Shylin26/swarmlens/internal/schema"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("expected a subcommand: replay")
	}

	switch os.Args[1] {
	case "replay":
		runReplay(os.Args[2:])
	default:
		log.Fatalf("unknown subcommand: %s", os.Args[1])
	}
}

func runReplay(args []string) {
	replayCmd := flag.NewFlagSet("replay", flag.ExitOnError)
	swarmID := replayCmd.String("swarm-id", "", "swarm ID to replay")
	from := replayCmd.String("from", "10m", "how far back to start replaying, e.g. 10m, 1h")
	replayCmd.Parse(args)

	if *swarmID == "" {
		log.Fatal("--swarm-id is required")
	}

	duration, err := time.ParseDuration(*from)
	if err != nil {
		log.Fatalf("invalid --from duration: %v", err)
	}
	fromTime := time.Now().Add(-duration)

	replayer := replay.NewReplayer([]string{"localhost:19092"}, "agent.messages")

	fmt.Printf("Replaying swarm %s from %s ago...\n\n", *swarmID, duration)

	err = replayer.Replay(context.Background(), *swarmID, fromTime, 6, func(event schema.Event) {
		content := ""
		if event.Payload.Content != nil {
			content = *event.Payload.Content
		}
		recipient := ""
		if event.Payload.RecipientAgentID != nil {
			recipient = *event.Payload.RecipientAgentID
		}
		fmt.Printf("[%s] %s -> %s: %s\n",
			event.Timestamp.Format("15:04:05"), event.AgentID, recipient, content)
	})
	if err != nil {
		log.Fatalf("replay failed: %v", err)
	}
}
