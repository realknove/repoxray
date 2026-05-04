package report

import (
	"encoding/json"
	"strings"
	"testing"

	"repoxray/internal/score"
	"repoxray/internal/types"
)

func TestRating(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{0, "Needs foundation"},
		{39, "Needs foundation"},
		{40, "Early-stage"},
		{59, "Early-stage"},
		{60, "Healthy"},
		{79, "Healthy"},
		{80, "Mature"},
		{89, "Mature"},
		{90, "Excellent"},
		{100, "Excellent"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			rating := Rating(tt.score)
			if rating != tt.expected {
				t.Errorf("Rating(%d) = %v, want %v", tt.score, rating, tt.expected)
			}
		})
	}
}

func TestRenderText(t *testing.T) {
	results := sampleResults()
	output, err := Render(results, score.Analyze(results), ".", TextFormat)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	wantParts := []string{
		"RepoXray Repository Health Report",
		"Total checks: 3",
		"Repository: .",
		"Passed: 1, Warned: 1, Failed: 1",
		"Maturity Score: 15/100",
		"Rating: Needs foundation",
		"Category Breakdown",
		"Documentation: 15/15 (100%) - excellent",
		"Licensing: 0/10 (0%) - at risk",
		"[PASS] README.md exists",
		"[WARN] No issue templates found",
		"[FAIL] No license file found",
		"High priority:\n- Add a LICENSE file.",
		"Medium priority:\n- Add issue templates.",
	}

	for _, part := range wantParts {
		if !strings.Contains(output, part) {
			t.Fatalf("text output missing %q:\n%s", part, output)
		}
	}
}

func TestRenderJSON(t *testing.T) {
	results := sampleResults()
	output, err := Render(results, score.Analyze(results), ".", JSONFormat)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	var got Report
	if err := json.Unmarshal([]byte(output), &got); err != nil {
		t.Fatalf("JSON output is not machine-readable: %v\n%s", err, output)
	}

	if got.Repository != "." {
		t.Fatalf("repository = %q, want %q", got.Repository, ".")
	}
	if got.Summary.Total != 3 || got.Summary.Passed != 1 || got.Summary.Warned != 1 || got.Summary.Failed != 1 {
		t.Fatalf("summary = %+v, want total 3 pass 1 warn 1 fail 1", got.Summary)
	}
	if got.Score.Maturity != 15 || got.Score.Max != 100 || got.Score.Rating != "Needs foundation" {
		t.Fatalf("score = %+v, want 15/100 Needs foundation", got.Score)
	}
	if len(got.Categories) != 7 {
		t.Fatalf("categories length = %d, want 7", len(got.Categories))
	}
	if got.Categories[0].ID != "documentation" || got.Categories[0].Percentage != 100 {
		t.Fatalf("first category = %+v, want documentation at 100 percent", got.Categories[0])
	}
	if len(got.Checks) != 3 {
		t.Fatalf("checks length = %d, want 3", len(got.Checks))
	}
	if len(got.Recommendations.High) != 1 || got.Recommendations.High[0] != "Add a LICENSE file." {
		t.Fatalf("high recommendations = %#v", got.Recommendations.High)
	}
}

func TestRenderMarkdown(t *testing.T) {
	results := sampleResults()
	output, err := Render(results, score.Analyze(results), "owner/repo", MarkdownFormat)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	wantParts := []string{
		"# RepoXray Repository Health Report",
		"## Summary",
		"| Repository | `owner/repo` |",
		"**Maturity Score:** 15/100",
		"## Category Breakdown",
		"| Documentation | 15/15 | 100% | excellent |",
		"### Category Recommendations",
		"## Strengths",
		"| PASS | README | README.md exists |",
		"## Recommendations",
		"### High Priority",
		"- Add a LICENSE file.",
		"### Medium Priority",
		"- Add issue templates.",
	}

	for _, part := range wantParts {
		if !strings.Contains(output, part) {
			t.Fatalf("markdown output missing %q:\n%s", part, output)
		}
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    Format
		wantErr bool
	}{
		{"text", TextFormat, false},
		{"json", JSONFormat, false},
		{"markdown", MarkdownFormat, false},
		{"JSON", JSONFormat, false},
		{"html", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseFormat(%q) returned nil error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseFormat(%q) returned error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func sampleResults() []types.Result {
	return []types.Result{
		{
			ID:          "readme",
			Title:       "README",
			Status:      types.Pass,
			Message:     "README.md exists",
			Points:      10,
			MaxPoints:   10,
			Description: "A README.md file provides essential information about the project.",
		},
		{
			ID:             "issue_templates",
			Title:          "Issue templates",
			Status:         types.Warn,
			Message:        "No issue templates found",
			Points:         0,
			MaxPoints:      5,
			Recommendation: "Add issue templates.",
		},
		{
			ID:             "license",
			Title:          "License",
			Status:         types.Fail,
			Message:        "No license file found",
			Points:         0,
			MaxPoints:      10,
			Recommendation: "Add a LICENSE file.",
		},
	}
}
