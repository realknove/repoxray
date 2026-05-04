package score

import (
	"repoxray/internal/types"
	"testing"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		results  []types.Result
		expected int
	}{
		{
			name:     "all pass",
			results:  allCategoryResults(true),
			expected: 100,
		},
		{
			name:     "all fail",
			results:  allCategoryResults(false),
			expected: 0,
		},
		{
			name: "category weights make documentation worth fifteen points",
			results: []types.Result{
				{ID: "readme", Points: 10, MaxPoints: 10},
			},
			expected: 15,
		},
		{
			name:     "empty results",
			results:  []types.Result{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateScore(tt.results)
			if score != tt.expected {
				t.Errorf("CalculateScore() = %v, want %v", score, tt.expected)
			}
		})
	}
}

func TestAnalyzeCategoryBreakdown(t *testing.T) {
	results := []types.Result{
		{ID: "readme", Points: 10, MaxPoints: 10},
		{ID: "license", Points: 0, MaxPoints: 10, Recommendation: "Add a LICENSE file."},
		{ID: "contributing", Points: 5, MaxPoints: 5},
		{ID: "code_of_conduct", Points: 0, MaxPoints: 5, Recommendation: "Add a CODE_OF_CONDUCT.md file."},
		{ID: "issue_templates", Points: 0, MaxPoints: 5, Recommendation: "Add issue templates."},
		{ID: "pr_template", Points: 5, MaxPoints: 5},
	}

	analysis := Analyze(results)

	if analysis.Overall != 22 {
		t.Fatalf("Overall = %d, want 22", analysis.Overall)
	}
	if analysis.Max != 100 {
		t.Fatalf("Max = %d, want 100", analysis.Max)
	}
	if len(analysis.Categories) != 7 {
		t.Fatalf("category count = %d, want 7", len(analysis.Categories))
	}

	documentation := findCategory(t, analysis, "documentation")
	if documentation.Score != 15 || documentation.MaxScore != 15 || documentation.Percentage != 100 || documentation.Status != "excellent" {
		t.Fatalf("documentation category = %+v, want 15/15 100 excellent", documentation)
	}

	licensing := findCategory(t, analysis, "licensing")
	if licensing.Score != 0 || licensing.MaxScore != 10 || licensing.Percentage != 0 || licensing.Status != "at risk" {
		t.Fatalf("licensing category = %+v, want 0/10 0 at risk", licensing)
	}
	if len(licensing.Recommendations) != 1 || licensing.Recommendations[0] != "Add a LICENSE file." {
		t.Fatalf("licensing recommendations = %#v", licensing.Recommendations)
	}

	contributor := findCategory(t, analysis, "contributor_readiness")
	if contributor.Score != 7 || contributor.MaxScore != 15 || contributor.Percentage != 50 || contributor.Status != "needs attention" {
		t.Fatalf("contributor category = %+v, want 7/15 50 needs attention", contributor)
	}
	if len(contributor.Recommendations) != 2 {
		t.Fatalf("contributor recommendations = %#v, want 2 recommendations", contributor.Recommendations)
	}
}

func TestCategoryStatus(t *testing.T) {
	tests := []struct {
		percentage int
		want       string
	}{
		{90, "excellent"},
		{89, "healthy"},
		{75, "healthy"},
		{74, "needs attention"},
		{50, "needs attention"},
		{49, "at risk"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := CategoryStatus(tt.percentage)
			if got != tt.want {
				t.Fatalf("CategoryStatus(%d) = %q, want %q", tt.percentage, got, tt.want)
			}
		})
	}
}

func findCategory(t *testing.T, analysis Analysis, id string) CategoryScore {
	t.Helper()
	for _, category := range analysis.Categories {
		if category.ID == id {
			return category
		}
	}
	t.Fatalf("category %q not found in %#v", id, analysis.Categories)
	return CategoryScore{}
}

func allCategoryResults(pass bool) []types.Result {
	points := map[string]int{
		"readme":                               10,
		"license":                              10,
		"contributing":                         5,
		"code_of_conduct":                      5,
		"issue_templates":                      5,
		"pr_template":                          5,
		"ci":                                   15,
		"tests":                                15,
		"security":                             5,
		"workflow_permissions":                 5,
		"workflow_write_all_permissions":       5,
		"workflow_action_pins":                 5,
		"workflow_pull_request_target":         5,
		"workflow_pull_request_target_secrets": 5,
		"git":                                  5,
		"package_files":                        10,
	}

	results := make([]types.Result, 0, len(points))
	for id, maxPoints := range points {
		resultPoints := 0
		if pass {
			resultPoints = maxPoints
		}
		results = append(results, types.Result{ID: id, Points: resultPoints, MaxPoints: maxPoints})
	}
	return results
}
