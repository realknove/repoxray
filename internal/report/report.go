package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/realknove/repoxray/internal/score"
	"github.com/realknove/repoxray/internal/types"
)

type Format string

const (
	TextFormat     Format = "text"
	JSONFormat     Format = "json"
	MarkdownFormat Format = "markdown"
)

type Report struct {
	Repository      string          `json:"repository"`
	Summary         Summary         `json:"summary"`
	Score           Score           `json:"score"`
	Categories      []Category      `json:"categories"`
	Checks          []types.Result  `json:"checks"`
	Strengths       []ReportItem    `json:"strengths"`
	Warnings        []ReportItem    `json:"warnings"`
	CriticalIssues  []ReportItem    `json:"critical_issues"`
	Recommendations Recommendations `json:"recommendations"`
}

type Summary struct {
	Total  int `json:"total"`
	Passed int `json:"passed"`
	Warned int `json:"warned"`
	Failed int `json:"failed"`
}

type Score struct {
	Maturity int    `json:"maturity"`
	Max      int    `json:"max"`
	Rating   string `json:"rating"`
}

type Category struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Score           int      `json:"score"`
	MaxScore        int      `json:"max_score"`
	Percentage      int      `json:"percentage"`
	Status          string   `json:"status"`
	Recommendations []string `json:"recommendations"`
}

type ReportItem struct {
	ID      string       `json:"id"`
	Title   string       `json:"title"`
	Status  types.Status `json:"status"`
	Message string       `json:"message"`
}

type Recommendations struct {
	High   []string `json:"high"`
	Medium []string `json:"medium"`
}

// Rating returns a maturity rating label based on the score
func Rating(score int) string {
	if score >= 90 {
		return "Excellent"
	} else if score >= 80 {
		return "Mature"
	} else if score >= 60 {
		return "Healthy"
	} else if score >= 40 {
		return "Early-stage"
	} else {
		return "Needs foundation"
	}
}

func ParseFormat(value string) (Format, error) {
	switch Format(strings.ToLower(strings.TrimSpace(value))) {
	case TextFormat:
		return TextFormat, nil
	case JSONFormat:
		return JSONFormat, nil
	case MarkdownFormat:
		return MarkdownFormat, nil
	default:
		return "", fmt.Errorf("unsupported format %q (expected text, json, or markdown)", value)
	}
}

func NewReport(results []types.Result, scoreAnalysis score.Analysis, repoPath string) Report {
	passed := 0
	warned := 0
	failed := 0
	strengths := []ReportItem{}
	warnings := []ReportItem{}
	criticalIssues := []ReportItem{}
	highRecs := []string{}
	mediumRecs := []string{}
	seenHighRecs := make(map[string]bool)
	seenMediumRecs := make(map[string]bool)

	for _, result := range results {
		switch result.Status {
		case types.Pass:
			passed++
			strengths = append(strengths, newReportItem(result))
		case types.Warn:
			warned++
			warnings = append(warnings, newReportItem(result))
			if result.Recommendation != "" && !seenMediumRecs[result.Recommendation] {
				mediumRecs = append(mediumRecs, result.Recommendation)
				seenMediumRecs[result.Recommendation] = true
			}
		case types.Fail:
			failed++
			criticalIssues = append(criticalIssues, newReportItem(result))
			if result.Recommendation != "" && !seenHighRecs[result.Recommendation] {
				highRecs = append(highRecs, result.Recommendation)
				seenHighRecs[result.Recommendation] = true
			}
		}
	}

	return Report{
		Repository: repoPath,
		Summary: Summary{
			Total:  len(results),
			Passed: passed,
			Warned: warned,
			Failed: failed,
		},
		Score: Score{
			Maturity: scoreAnalysis.Overall,
			Max:      scoreAnalysis.Max,
			Rating:   Rating(scoreAnalysis.Overall),
		},
		Categories:     newCategories(scoreAnalysis.Categories),
		Checks:         results,
		Strengths:      strengths,
		Warnings:       warnings,
		CriticalIssues: criticalIssues,
		Recommendations: Recommendations{
			High:   highRecs,
			Medium: mediumRecs,
		},
	}
}

func newCategories(categoryScores []score.CategoryScore) []Category {
	categories := make([]Category, 0, len(categoryScores))
	for _, categoryScore := range categoryScores {
		categories = append(categories, Category{
			ID:              categoryScore.ID,
			Name:            categoryScore.Name,
			Score:           categoryScore.Score,
			MaxScore:        categoryScore.MaxScore,
			Percentage:      categoryScore.Percentage,
			Status:          categoryScore.Status,
			Recommendations: categoryScore.Recommendations,
		})
	}
	return categories
}

func newReportItem(result types.Result) ReportItem {
	return ReportItem{
		ID:      result.ID,
		Title:   result.Title,
		Status:  result.Status,
		Message: result.Message,
	}
}

func Render(results []types.Result, scoreAnalysis score.Analysis, repoPath string, format Format) (string, error) {
	report := NewReport(results, scoreAnalysis, repoPath)

	switch format {
	case TextFormat:
		return renderText(report), nil
	case JSONFormat:
		return renderJSON(report)
	case MarkdownFormat:
		return renderMarkdown(report), nil
	default:
		return "", fmt.Errorf("unsupported format %q", format)
	}
}

func renderJSON(report Report) (string, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func renderText(report Report) string {
	var b strings.Builder

	b.WriteString("RepoXray Repository Health Report\n")
	b.WriteString("=================================\n\n")

	b.WriteString("Summary\n")
	b.WriteString("-------\n")
	fmt.Fprintf(&b, "Total checks: %d\n", report.Summary.Total)
	fmt.Fprintf(&b, "Repository: %s\n", report.Repository)
	fmt.Fprintf(&b, "Passed: %d, Warned: %d, Failed: %d\n", report.Summary.Passed, report.Summary.Warned, report.Summary.Failed)
	b.WriteString("\n")

	b.WriteString("Score\n")
	b.WriteString("-----\n")
	fmt.Fprintf(&b, "Maturity Score: %d/%d\n", report.Score.Maturity, report.Score.Max)
	fmt.Fprintf(&b, "Rating: %s\n", report.Score.Rating)
	b.WriteString("\n")

	b.WriteString("Category Breakdown\n")
	b.WriteString("------------------\n")
	for _, category := range report.Categories {
		fmt.Fprintf(&b, "%s: %d/%d (%d%%) - %s\n", category.Name, category.Score, category.MaxScore, category.Percentage, category.Status)
		if len(category.Recommendations) > 0 {
			b.WriteString("  Recommendations:\n")
			for _, rec := range category.Recommendations {
				fmt.Fprintf(&b, "  - %s\n", rec)
			}
		}
	}
	b.WriteString("\n")

	b.WriteString("Strengths\n")
	b.WriteString("---------\n")
	if len(report.Strengths) > 0 {
		for _, item := range report.Strengths {
			fmt.Fprintf(&b, "[PASS] %s\n", item.Message)
		}
	} else {
		b.WriteString("No strengths identified.\n")
	}
	b.WriteString("\n")

	b.WriteString("Warnings\n")
	b.WriteString("--------\n")
	if len(report.Warnings) > 0 {
		for _, item := range report.Warnings {
			fmt.Fprintf(&b, "[WARN] %s\n", item.Message)
		}
	} else {
		b.WriteString("No warnings.\n")
	}
	b.WriteString("\n")

	b.WriteString("Critical Issues\n")
	b.WriteString("---------------\n")
	if len(report.CriticalIssues) > 0 {
		for _, item := range report.CriticalIssues {
			fmt.Fprintf(&b, "[FAIL] %s\n", item.Message)
		}
	} else {
		b.WriteString("No critical issues found.\n")
	}
	b.WriteString("\n")

	b.WriteString("Recommendations\n")
	b.WriteString("---------------\n")
	hasRecs := false
	if len(report.Recommendations.High) > 0 {
		b.WriteString("High priority:\n")
		for _, rec := range report.Recommendations.High {
			fmt.Fprintf(&b, "- %s\n", rec)
		}
		hasRecs = true
	}
	if len(report.Recommendations.Medium) > 0 {
		if hasRecs {
			b.WriteString("\n")
		}
		b.WriteString("Medium priority:\n")
		for _, rec := range report.Recommendations.Medium {
			fmt.Fprintf(&b, "- %s\n", rec)
		}
		hasRecs = true
	}
	if !hasRecs {
		b.WriteString("No recommendations. Repository looks healthy.\n")
	}

	return b.String()
}

func renderMarkdown(report Report) string {
	var b strings.Builder

	b.WriteString("# RepoXray Repository Health Report\n\n")
	b.WriteString("## Summary\n\n")
	b.WriteString("| Metric | Value |\n")
	b.WriteString("| --- | --- |\n")
	fmt.Fprintf(&b, "| Repository | `%s` |\n", escapeMarkdownTable(report.Repository))
	fmt.Fprintf(&b, "| Total checks | %d |\n", report.Summary.Total)
	fmt.Fprintf(&b, "| Passed | %d |\n", report.Summary.Passed)
	fmt.Fprintf(&b, "| Warned | %d |\n", report.Summary.Warned)
	fmt.Fprintf(&b, "| Failed | %d |\n\n", report.Summary.Failed)

	b.WriteString("## Score\n\n")
	fmt.Fprintf(&b, "**Maturity Score:** %d/%d\n\n", report.Score.Maturity, report.Score.Max)
	fmt.Fprintf(&b, "**Rating:** %s\n\n", report.Score.Rating)

	b.WriteString("## Category Breakdown\n\n")
	b.WriteString("| Category | Score | Percentage | Status |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	for _, category := range report.Categories {
		fmt.Fprintf(&b, "| %s | %d/%d | %d%% | %s |\n", escapeMarkdownTable(category.Name), category.Score, category.MaxScore, category.Percentage, escapeMarkdownTable(category.Status))
	}
	b.WriteString("\n")

	hasCategoryRecs := false
	for _, category := range report.Categories {
		if len(category.Recommendations) == 0 {
			continue
		}
		if !hasCategoryRecs {
			b.WriteString("### Category Recommendations\n\n")
			hasCategoryRecs = true
		}
		fmt.Fprintf(&b, "**%s**\n\n", escapeMarkdownTable(category.Name))
		for _, rec := range category.Recommendations {
			fmt.Fprintf(&b, "- %s\n", rec)
		}
		b.WriteString("\n")
	}

	writeMarkdownItems(&b, "Strengths", report.Strengths, "No strengths identified.")
	writeMarkdownItems(&b, "Warnings", report.Warnings, "No warnings.")
	writeMarkdownItems(&b, "Critical Issues", report.CriticalIssues, "No critical issues found.")

	b.WriteString("## Recommendations\n\n")
	if len(report.Recommendations.High) == 0 && len(report.Recommendations.Medium) == 0 {
		b.WriteString("No recommendations. Repository looks healthy.\n")
		return b.String()
	}
	if len(report.Recommendations.High) > 0 {
		b.WriteString("### High Priority\n\n")
		for _, rec := range report.Recommendations.High {
			fmt.Fprintf(&b, "- %s\n", rec)
		}
		b.WriteString("\n")
	}
	if len(report.Recommendations.Medium) > 0 {
		b.WriteString("### Medium Priority\n\n")
		for _, rec := range report.Recommendations.Medium {
			fmt.Fprintf(&b, "- %s\n", rec)
		}
	}

	return b.String()
}

func writeMarkdownItems(b *strings.Builder, title string, items []ReportItem, emptyMessage string) {
	fmt.Fprintf(b, "## %s\n\n", title)
	if len(items) == 0 {
		fmt.Fprintf(b, "%s\n\n", emptyMessage)
		return
	}

	b.WriteString("| Status | Check | Message |\n")
	b.WriteString("| --- | --- | --- |\n")
	for _, item := range items {
		fmt.Fprintf(
			b,
			"| %s | %s | %s |\n",
			strings.ToUpper(string(item.Status)),
			escapeMarkdownTable(item.Title),
			escapeMarkdownTable(item.Message),
		)
	}
	b.WriteString("\n")
}

func escapeMarkdownTable(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "|", `\|`)
	value = strings.ReplaceAll(value, "\n", " ")
	return value
}
