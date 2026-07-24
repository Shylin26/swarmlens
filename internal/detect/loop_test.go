package detect

import "testing"

func TestLoopDetected_TrueOnRepeatingPattern(t *testing.T) {
	window := []string{
		"worker->reviewer", "reviewer->worker",
		"worker->reviewer", "reviewer->worker",
		"worker->reviewer", "reviewer->worker",
	}
	if !LoopDetected(window, 3) {
		t.Error("expected loop to be detected ,got false")
	}
}
func TestLoopDetected_FalseOnHealthySwarm(t *testing.T) {
	window := []string{
		"reviewer->planner",
		"worker->reviewer",
		"planner->worker",
	}
	if LoopDetected(window, 2) {
		t.Error("expected no loop, got true")
	}
}
func TestLoopDetected_FalseWhenWindwoooShort(t *testing.T) {
	window := []string{"worker->reviewer", "reviewer->worker"}
	if LoopDetected(window, 3) {
		t.Error("expected no loop when window shorter than required length, got true")
	}
}
