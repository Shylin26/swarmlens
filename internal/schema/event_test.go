package schema

import "testing"

func TestNewEvent_IsValid(t *testing.T) {
	meta := Metadata{Framework: "custom", SDKVersion: "0.1.0"}
	event := NewEvent("swarm-1", "agent-1", EventMessage, Payload{}, meta)
	if err := event.Validate(); err != nil {
		t.Errorf("expected valid event,got error:%v", err)
	}
}
func TestEvent_MissingAgentID_FailsValidation(t *testing.T) {
	meta := Metadata{Framework: "custom", SDKVersion: "0.1.0"}
	event := NewEvent("swarm-1", "", EventMessage, Payload{}, meta)
	if err := event.Validate(); err == nil {
		t.Error("expected validation error for missing agent_id,got nil")
	}
}
