package main

import (
	"context"
	"log"

	"github.com/Shylin26/swarmlens/internal/state"
)

func main() {
	ctx := context.Background()
	store := state.NewStore("localhost:16379")
	if err := store.RecordAgentActivity(ctx, "swarm-demo-1", "planner"); err != nil {
		log.Fatalf("failed to record activity :%v", err)
	}
	log.Println("recorded agent activity successfully")
}
