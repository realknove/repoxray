package score

import (
	"repoxray/internal/types"
)

// CalculateScore computes the maturity score from check results
func CalculateScore(results []types.Result) int {
	totalPoints := 0
	totalMaxPoints := 0

	for _, result := range results {
		totalPoints += result.Points
		totalMaxPoints += result.MaxPoints
	}

	if totalMaxPoints == 0 {
		return 0
	}

	score := (totalPoints * 100) / totalMaxPoints
	return score
}
