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
			name: "all pass",
			results: []types.Result{
				{Points: 10, MaxPoints: 10},
				{Points: 5, MaxPoints: 5},
			},
			expected: 100,
		},
		{
			name: "half pass",
			results: []types.Result{
				{Points: 10, MaxPoints: 10},
				{Points: 0, MaxPoints: 10},
			},
			expected: 50,
		},
		{
			name: "no points",
			results: []types.Result{
				{Points: 0, MaxPoints: 10},
				{Points: 0, MaxPoints: 5},
			},
			expected: 0,
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
