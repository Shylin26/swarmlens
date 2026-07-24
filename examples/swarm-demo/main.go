package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Shylin26/swarmlens/internal/ingest"
	"github.com/Shylin26/swarmlens/internal/llm"
)

func main() {
	ctx := context.Background()
	producer := ingest.NewProducer([]string{"localhost:19092"}, "agent.messages")
	model := llm.NewOllamaClient("http://localhost:11434", "mistral:7b")
	planner := NewAgent("planner", producer, model)
	worker := NewAgent("worker", producer, model)
	reviewer := NewAgent("reviewer", producer, model)

	swarmID := "swarm-demo-1"

	log.Println("planner:breaking down the task...")
	instruction, err := planner.Act(ctx, swarmID, worker.ID, nil,
		"You are a planner agent. Give the worker agent a single clear one-sentence instruction to write a 3-line product description for a ceramic coffee mug. Reply with only the instruction, nothing else.")
	if err != nil {
		log.Fatalf("planner failed :%v", err)
	}
	fmt.Printf("\n[planner -> worker]: %s\n", instruction)
	taskID := "task-1"
	log.Println("worker:writing the description...")
	draft, err := worker.Act(ctx, swarmID, reviewer.ID, &taskID,
		"You are a worker agent. Follow this instruction exactly and reply with only the result: "+instruction)
	if err != nil {
		log.Fatalf("worker failed: %v", err)
	}
	fmt.Printf("\n[worker -> reviewer]: %s\n", draft)
	log.Println("reviewer:judging the draft...")
	verdict, err := reviewer.Act(ctx, swarmID, planner.ID, &taskID,
		"You are a reviewer agent. Judge this product description draft. Reply with either 'APPROVED' or 'NEEDS REVISION' followed by a one-sentence reason. Draft: "+draft)
	if err != nil {
		log.Fatalf("reviewer failed: %v", err)

	}
	fmt.Printf("\n[reviewer -> planner]: %s\n", verdict)
	log.Println("forcing an artificial loop between worker and reviewer...")
	for i := 0; i < 4; i++ {
		if _, err := worker.Act(ctx, swarmID, reviewer.ID, &taskID, "ping"); err != nil {
			log.Fatalf("worker failed during forced loop: %v", err)
		}
		if _, err := reviewer.Act(ctx, swarmID, worker.ID, &taskID, "pong"); err != nil {
			log.Fatalf("reviewer failed during forced loop: %v", err)
		}
	}

	log.Println("swarm run complete")

}
