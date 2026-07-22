package schema

import (
	"time"

	"errors"

	"github.com/google/uuid"
)

type EventType string

const (
	EventMessage    EventType = "message"
	EventToolCall   EventType = "tool_call"
	EventToolResult EventType = "tool_result"
	EventLifeCycle  EventType = "lifecycle"
	EventDecision   EventType = "decision"
)

type Payload struct {
	Content          *string        `json:"content,omitempty"`
	RecipientAgentID *string        `json:"recipient_agent_id,omitempty"`
	ToolName         *string        `json:"tool_name,omitempty"`
	ToolArgs         map[string]any `json:"tool_args,omitempty"`
	TokensIn         *int           `json:"tokens_in,omitempty"`
	TokensOut        *int           `json:"tokens_out,omitempty"`
	CostUSD          *float64       `json:"cost_usd,omitempty"`
	Model            *string        `json:"model,omitempty"`
}

type Event struct {
	EventID      string    `json:"event_id"`
	SwarmID      string    `json:"swarm_id"`
	AgentID      string    `json:"agent_id"`
	ParentTaskID *string   `json"parent_task_id",omniempty"`
	Type         EventType `json:"event_type"`
	Timestamp    time.Time `json:"timestamp"`
	Payload      Payload   `json:"payload"`
	Metadata     Metadata  `json:"metadata"`
}
type Metadata struct {
	Framework  string `json:"framework"`
	SDKVersion string `json:"sdk_version"`
}

func NewEvent(swarmID, agentID string, eventType EventType, payload Payload, meta Metadata) Event {
	return Event{
		EventID:   uuid.New().String(),
		SwarmID:   swarmID,
		AgentID:   agentID,
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
		Metadata:  meta,
	}
}
func (e Event) Validate() error {
	if e.EventID == "" {
		return errors.New("event :missing event_id")

	}
	if e.SwarmID == "" {
		return errors.New("event:missing swarm_id")
	}
	if e.AgentID == "" {
		return errors.New("event:missing agent_id")
	}
	if e.Type == "" {
		return errors.New("event:missing event_type")
	}
	if e.AgentID == "" {
		return errors.New("event:missing timestamp")
	}
	return nil

}
