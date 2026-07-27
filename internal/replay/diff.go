package replay

import (
	"github.com/Shylin26/swarmlens/internal/detect"
	"github.com/Shylin26/swarmlens/internal/schema"
)

// DiffResult summarizes a comparison between two swarms' event histories.
type DiffResult struct {
	SwarmAAlerts []string
	SwarmBAlerts []string
	SwarmARoles  detect.RoleCollapseResult
	SwarmBRoles  detect.RoleCollapseResult
}

// AnalyzeSwarm runs the loop and role-collapse detectors over a full
// event history and returns any alerts found, plus the final role distribution.
func AnalyzeSwarm(events []schema.Event) (alerts []string, roles detect.RoleCollapseResult) {
	var window []string
	var completions []string

	for _, event := range events {
		if event.Type != schema.EventMessage {
			continue
		}

		if event.Payload.RecipientAgentID != nil {
			entry := event.AgentID + "->" + *event.Payload.RecipientAgentID
			window = append([]string{entry}, window...)
			if len(window) > 20 {
				window = window[:20]
			}
			if detect.LoopDetected(window, 3) {
				alerts = append(alerts, "loop detected involving "+event.AgentID)
			}
		}

		if event.ParentTaskID != nil {
			completions = append(completions, event.AgentID)
		}
	}

	roles = detect.DetectRoleCollapse(completions, 0.7)
	if roles.Collapsed {
		alerts = append(alerts, "role collapse: "+roles.DominantAgent)
	}

	return alerts, roles
}
