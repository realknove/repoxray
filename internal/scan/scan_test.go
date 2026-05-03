package scan

import (
	"os"
	"path/filepath"
	"testing"

	"repoxray/internal/types"
)

type stubCheck struct {
	result types.Result
	seen   *types.Context
}

func (c stubCheck) Run(ctx types.Context) types.Result {
	if c.seen != nil {
		*c.seen = ctx
	}
	return c.result
}

func TestScan(t *testing.T) {
	// Create temp dir
	tempDir, err := os.MkdirTemp("", "repoxray_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create a README.md
	readmePath := filepath.Join(tempDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Scan
	results := Scan(tempDir)

	// Check that we have results
	if len(results) == 0 {
		t.Errorf("Expected at least one result, got 0")
	}

	// Check that README check passed
	found := false
	for _, result := range results {
		if result.ID == "readme" {
			found = true
			if result.Status != types.Pass {
				t.Errorf("Expected README to pass, got %v", result.Status)
			}
			break
		}
	}
	if !found {
		t.Errorf("README check not found in results")
	}

	// Check total points
	totalPoints := 0
	totalMaxPoints := 0
	for _, result := range results {
		totalPoints += result.Points
		totalMaxPoints += result.MaxPoints
	}
	if totalMaxPoints == 0 {
		t.Errorf("Total max points should not be 0")
	}
	score := (totalPoints * 100) / totalMaxPoints
	if score < 0 || score > 100 {
		t.Errorf("Score out of range: %d", score)
	}
}

func TestScanWithChecksAggregatesResults(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "repoxray_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.WriteFile(filepath.Join(tempDir, "README.md"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	testDir := filepath.Join(tempDir, "internal", "example")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testDir, "example_test.go"), []byte("package example"), 0644); err != nil {
		t.Fatal(err)
	}

	var seen types.Context
	results := ScanWithChecks(tempDir, []types.Check{
		stubCheck{
			result: types.Result{ID: "first", Status: types.Pass, Points: 1, MaxPoints: 1},
			seen:   &seen,
		},
		stubCheck{
			result: types.Result{ID: "second", Status: types.Warn, Points: 0, MaxPoints: 1},
		},
	})

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	if results[0].ID != "first" || results[1].ID != "second" {
		t.Fatalf("Expected scanner to preserve check order, got %q then %q", results[0].ID, results[1].ID)
	}
	if seen.RepoPath != tempDir {
		t.Fatalf("Expected context repo path %q, got %q", tempDir, seen.RepoPath)
	}
	if !seen.HasPath("README.md") {
		t.Fatal("Expected context to include discovered README.md")
	}
	if !seen.HasFileWithSuffix("_test.go") {
		t.Fatal("Expected context to include discovered test files")
	}
}
