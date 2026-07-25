package detect

type RoleCollapseResult struct {
	Collapsed     bool
	DominantAgent string
	Share         float64
}

func DetectRoleCollapse(completions []string, threshold float64) RoleCollapseResult {
	if len(completions) == 0 {
		return RoleCollapseResult{}
	}
	counts := make(map[string]int)
	for _, agentID := range completions {
		counts[agentID]++
	}
	var dominantAgent string
	var dominantCount int
	for agentID, count := range counts {
		if count > dominantCount {
			dominantAgent = agentID
			dominantCount = count
		}
	}
	share := float64(dominantCount) / float64(len(completions))
	return RoleCollapseResult{
		Collapsed:     share > threshold,
		DominantAgent: dominantAgent,
		Share:         share,
	}
}
