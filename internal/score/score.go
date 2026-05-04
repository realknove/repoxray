package score

import (
	"repoxray/internal/types"
)

type Analysis struct {
	Overall    int             `json:"overall"`
	Max        int             `json:"max"`
	Categories []CategoryScore `json:"categories"`
}

type CategoryScore struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Score           int      `json:"score"`
	MaxScore        int      `json:"max_score"`
	Percentage      int      `json:"percentage"`
	Status          string   `json:"status"`
	Recommendations []string `json:"recommendations"`
}

type categoryDefinition struct {
	ID       string
	Name     string
	MaxScore int
	CheckIDs []string
}

var categories = []categoryDefinition{
	{
		ID:       "documentation",
		Name:     "Documentation",
		MaxScore: 15,
		CheckIDs: []string{"readme"},
	},
	{
		ID:       "licensing",
		Name:     "Licensing",
		MaxScore: 10,
		CheckIDs: []string{"license"},
	},
	{
		ID:       "contributor_readiness",
		Name:     "Contributor readiness",
		MaxScore: 15,
		CheckIDs: []string{"contributing", "code_of_conduct", "issue_templates", "pr_template"},
	},
	{
		ID:       "ci_and_automation",
		Name:     "CI and automation",
		MaxScore: 20,
		CheckIDs: []string{"ci", "tests"},
	},
	{
		ID:       "security_posture",
		Name:     "Security posture",
		MaxScore: 20,
		CheckIDs: []string{"security", "workflow_permissions", "workflow_write_all_permissions", "workflow_action_pins", "workflow_pull_request_target", "workflow_pull_request_target_secrets"},
	},
	{
		ID:       "maintenance_signals",
		Name:     "Maintenance signals",
		MaxScore: 10,
		CheckIDs: []string{"git"},
	},
	{
		ID:       "project_structure",
		Name:     "Project structure",
		MaxScore: 10,
		CheckIDs: []string{"package_files"},
	},
}

// CalculateScore computes the weighted maturity score from check results.
func CalculateScore(results []types.Result) int {
	return Analyze(results).Overall
}

func Analyze(results []types.Result) Analysis {
	resultsByID := make(map[string]types.Result, len(results))
	for _, result := range results {
		resultsByID[result.ID] = result
	}

	categoryScores := make([]CategoryScore, 0, len(categories))
	overall := 0
	max := 0

	for _, category := range categories {
		categoryScore := calculateCategoryScore(category, resultsByID)
		categoryScores = append(categoryScores, categoryScore)
		overall += categoryScore.Score
		max += categoryScore.MaxScore
	}

	return Analysis{
		Overall:    overall,
		Max:        max,
		Categories: categoryScores,
	}
}

func calculateCategoryScore(category categoryDefinition, resultsByID map[string]types.Result) CategoryScore {
	points := 0
	maxPoints := 0
	recommendations := []string{}
	seenRecommendations := make(map[string]bool)

	for _, checkID := range category.CheckIDs {
		result, ok := resultsByID[checkID]
		if !ok {
			continue
		}

		points += result.Points
		maxPoints += result.MaxPoints
		if result.Recommendation != "" && !seenRecommendations[result.Recommendation] {
			recommendations = append(recommendations, result.Recommendation)
			seenRecommendations[result.Recommendation] = true
		}
	}

	percentage := 0
	if maxPoints > 0 {
		percentage = (points * 100) / maxPoints
	}
	weightedScore := (percentage * category.MaxScore) / 100

	return CategoryScore{
		ID:              category.ID,
		Name:            category.Name,
		Score:           weightedScore,
		MaxScore:        category.MaxScore,
		Percentage:      percentage,
		Status:          CategoryStatus(percentage),
		Recommendations: recommendations,
	}
}

func CategoryStatus(percentage int) string {
	if percentage >= 90 {
		return "excellent"
	}
	if percentage >= 75 {
		return "healthy"
	}
	if percentage >= 50 {
		return "needs attention"
	}
	return "at risk"
}
