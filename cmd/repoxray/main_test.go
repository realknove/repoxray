package main

import "testing"

func TestParseGitHubRepo(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOwner  string
		wantRepo   string
		wantParsed bool
	}{
		{
			name:       "owner slash repo",
			input:      "biomejs/biome",
			wantOwner:  "biomejs",
			wantRepo:   "biome",
			wantParsed: true,
		},
		{
			name:       "github host",
			input:      "github.com/biomejs/biome",
			wantOwner:  "biomejs",
			wantRepo:   "biome",
			wantParsed: true,
		},
		{
			name:       "https github URL",
			input:      "https://github.com/biomejs/biome",
			wantOwner:  "biomejs",
			wantRepo:   "biome",
			wantParsed: true,
		},
		{
			name:       "repo git suffix",
			input:      "github.com/biomejs/biome.git",
			wantOwner:  "biomejs",
			wantRepo:   "biome",
			wantParsed: true,
		},
		{
			name:       "trailing slash",
			input:      "github.com/biomejs/biome/",
			wantOwner:  "biomejs",
			wantRepo:   "biome",
			wantParsed: true,
		},
		{
			name:       "local current directory",
			input:      ".",
			wantParsed: false,
		},
		{
			name:       "absolute local path",
			input:      "/tmp/repo",
			wantParsed: false,
		},
		{
			name:       "too many path segments",
			input:      "github.com/biomejs/biome/issues",
			wantParsed: false,
		},
		{
			name:       "unsupported host",
			input:      "gitlab.com/biomejs/biome",
			wantParsed: false,
		},
		{
			name:       "empty owner",
			input:      "/biome",
			wantParsed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseGitHubRepo(tt.input)
			if ok != tt.wantParsed {
				t.Fatalf("parseGitHubRepo(%q) parsed = %v, want %v", tt.input, ok, tt.wantParsed)
			}
			if !tt.wantParsed {
				return
			}
			if got.Owner != tt.wantOwner || got.Name != tt.wantRepo {
				t.Fatalf("parseGitHubRepo(%q) = %s/%s, want %s/%s", tt.input, got.Owner, got.Name, tt.wantOwner, tt.wantRepo)
			}
		})
	}
}
