package main

import (
	"context"
	"log"

	"github.com/Shylin26/swarmlens/internal/llm"
)

func main() {
	client := llm.NewOllamaClient("http://localhost:11434", "mistral:7b")
	response, err := client.Generate(context.Background(), "Say hello in exactly in five words")
	if err != nil {
		log.Fatalf("failed to generate: %v", err)
	}
	log.Printf("model responded: %s", response)
}
