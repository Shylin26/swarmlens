package main

import (
	"context"
	"fmt"

	"github.com/Shylin26/swarmlens/internal/ingest"
	"github.com/Shylin26/swarmlens/internal/llm"
	"github.com/Shylin26/swarmlens/internal/schema"
)

type Agent struct {
	ID       string
	producer *ingest.Producer
	model    *llm.OllamaClient
}

func NewAgent(id string, producer *ingest.Producer, model *llm.OllamaClient) *Agent {
	return &Agent{ID: id, producer: producer, model: model}
}

func (a *Agent) Act(ctx context.Context, swarmID string, recipientID string, parentTaskID *string, prompt string) (string, error) {
	response, err := a.model.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("agent %s: failed to generate: %w", a.ID, err)
	}

	tokensOut := len(response) / 4
	costUSD := float64(tokensOut) * 0.000002
	meta := schema.Metadata{Framework: "custom", SDKVersion: "0.1.0"}
	event := schema.NewEvent(swarmID, a.ID, schema.EventMessage, schema.Payload{
		Content:          &response,
		RecipientAgentID: &recipientID,
		TokensOut:        &tokensOut,
		CostUSD:          &costUSD,
	}, meta)
	event.ParentTaskID = parentTaskID

	if err := a.producer.Publish(ctx, event); err != nil {
		return "", fmt.Errorf("agent %s: failed to publish: %w", a.ID, err)
	}

	return response, nil
}
