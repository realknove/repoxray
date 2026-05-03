package report

import (
	"testing"
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
