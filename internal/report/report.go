package report

import (
	"fmt"
	"repoxray/internal/types"
)

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

// RenderReport prints the diagnostic report
func RenderReport(results []types.Result, maturityScore int, repoPath string) {
	fmt.Println("RepoXray Repository Health Report")
	fmt.Println("=================================")
	fmt.Println()

	// Summary
	fmt.Println("Summary")
	fmt.Println("-------")
	passed := 0
	warned := 0
	failed := 0
	for _, result := range results {
		switch result.Status {
		case types.Pass:
			passed++
		case types.Warn:
			warned++
		case types.Fail:
			failed++
		}
	}
	fmt.Printf("Total checks: %d\n", len(results))
	fmt.Printf("Repository: %s\n", repoPath)
	fmt.Printf("Passed: %d, Warned: %d, Failed: %d\n", passed, warned, failed)
	fmt.Println()

	// Score
	fmt.Println("Score")
	fmt.Println("-----")
	fmt.Printf("Maturity Score: %d/100\n", maturityScore)
	fmt.Printf("Rating: %s\n", Rating(maturityScore))
	fmt.Println()

	// Strengths
	fmt.Println("Strengths")
	fmt.Println("---------")
	strengths := []string{}
	for _, result := range results {
		if result.Status == types.Pass {
			strengths = append(strengths, result.Message)
		}
	}
	if len(strengths) > 0 {
		for _, s := range strengths {
			fmt.Printf("[PASS] %s\n", s)
		}
	} else {
		fmt.Println("No strengths identified.")
	}
	fmt.Println()

	// Warnings
	fmt.Println("Warnings")
	fmt.Println("--------")
	warnings := []string{}
	for _, result := range results {
		if result.Status == types.Warn {
			warnings = append(warnings, result.Message)
		}
	}
	if len(warnings) > 0 {
		for _, w := range warnings {
			fmt.Printf("[WARN] %s\n", w)
		}
	} else {
		fmt.Println("No warnings.")
	}
	fmt.Println()

	// Critical Issues
	fmt.Println("Critical Issues")
	fmt.Println("---------------")
	missing := []string{}
	for _, result := range results {
		if result.Status == types.Fail {
			missing = append(missing, result.Message)
		}
	}
	if len(missing) > 0 {
		for _, m := range missing {
			fmt.Printf("[FAIL] %s\n", m)
		}
	} else {
		fmt.Println("No critical issues found.")
	}
	fmt.Println()

	// Recommendations
	fmt.Println("Recommendations")
	fmt.Println("---------------")
	highRecs := make(map[string]bool)
	mediumRecs := make(map[string]bool)
	for _, result := range results {
		if result.Recommendation != "" {
			if result.Status == types.Fail {
				highRecs[result.Recommendation] = true
			} else if result.Status == types.Warn {
				mediumRecs[result.Recommendation] = true
			}
		}
	}
	hasRecs := false
	if len(highRecs) > 0 {
		fmt.Println("High priority:")
		for rec := range highRecs {
			fmt.Printf("- %s\n", rec)
		}
		hasRecs = true
	}
	if len(mediumRecs) > 0 {
		if hasRecs {
			fmt.Println()
		}
		fmt.Println("Medium priority:")
		for rec := range mediumRecs {
			fmt.Printf("- %s\n", rec)
		}
		hasRecs = true
	}
	if !hasRecs {
		fmt.Println("No recommendations. Repository looks healthy.")
	}
}
