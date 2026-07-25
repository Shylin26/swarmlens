package detect

import "testing"

func TestDetectRoleCollapse_FlagsADominantAgent(t *testing.T) {
	completions := []string{
		"worker", "worker", "worker", "worker", "worker", "worker", "worker", "worker",
		"planner", "reviewer",
	}
	result := DetectRoleCollapse(completions, 0.7)

	if !result.Collapsed {
		t.Error("expected role collapse to be detected, got false")
	}
	if result.DominantAgent != "worker" {
		t.Errorf("expected dominant agent to be worker, got %s", result.DominantAgent)
	}

}

func TestDetectRoleCollapse_BalancedSwarmNotFlagged(t *testing.T) {
	completions := []string{
		"planner", "worker", "reviewer",
		"planner", "worker", "reviewer",
		"planner", "worker", "reviewer",
	}
	result := DetectRoleCollapse(completions, 0.7)
	if result.Collapsed {
		t.Errorf("expected balanced swarm to not be flagged, got collapsed with dominant agent %s at share %v",
			result.DominantAgent, result.Share)
	}

}
func TestDetectRoleCollapse_EmptyInputNotFlagged(t *testing.T) {
	result := DetectRoleCollapse([]string{}, 0.7)
	if result.Collapsed {
		t.Error("expected empty completions to not be flagged, got true")
	}
}

func TestDetectRoleCollapse_ExactlyAtThresholdNotFlagged(t *testing.T) {
	completions := []string{"worker", "worker", "worker", "worker", "worker", "worker", "worker", "planner", "reviewer", "designer"}
	result := DetectRoleCollapse(completions, 0.7)
	if result.Collapsed {
		t.Error("expected exactly-at-threshold share to not be flagged (boundary is exclusive), got true")
	}

}
