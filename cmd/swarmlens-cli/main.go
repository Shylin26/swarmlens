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
	case "diff":
		runDiff(os.Args[2:])
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
func runDiff(args []string) {
	diffCmd := flag.NewFlagSet("diff", flag.ExitOnError)
	swarmA := diffCmd.String("swarm-a", "", "first swarm ID to compare")
	swarmB := diffCmd.String("swarm-b", "", "second swarm ID to compare")
	from := diffCmd.String("from", "72h", "how far back to look for both swarms")
	diffCmd.Parse(args)

	if *swarmA == "" || *swarmB == "" {
		log.Fatal("--swarm-a and --swarm-b are both required")
	}

	duration, err := time.ParseDuration(*from)
	if err != nil {
		log.Fatalf("invalid --from duration: %v", err)
	}
	fromTime := time.Now().Add(-duration)

	replayer := replay.NewReplayer([]string{"localhost:19092"}, "agent.messages")
	ctx := context.Background()

	eventsA, err := replayer.CollectEvents(ctx, *swarmA, fromTime, 6)
	if err != nil {
		log.Fatalf("failed to collect events for swarm A: %v", err)
	}
	eventsB, err := replayer.CollectEvents(ctx, *swarmB, fromTime, 6)
	if err != nil {
		log.Fatalf("failed to collect events for swarm B: %v", err)
	}

	alertsA, rolesA := replay.AnalyzeSwarm(eventsA)
	alertsB, rolesB := replay.AnalyzeSwarm(eventsB)

	fmt.Printf("Swarm A (%s): %d events, %d alerts\n", *swarmA, len(eventsA), len(alertsA))
	for _, alert := range alertsA {
		fmt.Printf("  - %s\n", alert)
	}
	fmt.Printf("  role distribution: %s at %.0f%%\n\n", rolesA.DominantAgent, rolesA.Share*100)

	fmt.Printf("Swarm B (%s): %d events, %d alerts\n", *swarmB, len(eventsB), len(alertsB))
	for _, alert := range alertsB {
		fmt.Printf("  - %s\n", alert)
	}
	fmt.Printf("  role distribution: %s at %.0f%%\n", rolesB.DominantAgent, rolesB.Share*100)
}
