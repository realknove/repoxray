package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"repoxray/internal/report"
	"repoxray/internal/scan"
	"repoxray/internal/score"
)

const version = "0.1.0"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printHelp()
		return
	}

	command := args[0]

	switch command {
	case "scan":
		repoPath, outputFormat, err := parseScanArgs(args[1:])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		scanRepo(repoPath, outputFormat)
	case "version":
		fmt.Printf("RepoXray version %s\n", version)
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func parseScanArgs(args []string) (string, report.Format, error) {
	outputFormat := report.TextFormat
	repoPath := ""

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--format":
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("--format requires a value")
			}
			parsed, err := report.ParseFormat(args[i+1])
			if err != nil {
				return "", "", err
			}
			outputFormat = parsed
			i++
		case strings.HasPrefix(arg, "--format="):
			parsed, err := report.ParseFormat(strings.TrimPrefix(arg, "--format="))
			if err != nil {
				return "", "", err
			}
			outputFormat = parsed
		case strings.HasPrefix(arg, "-"):
			return "", "", fmt.Errorf("unknown scan option %q", arg)
		case repoPath == "":
			repoPath = arg
		default:
			return "", "", fmt.Errorf("scan accepts only one path")
		}
	}

	if repoPath == "" {
		return "", "", fmt.Errorf("scan command requires a path")
	}

	return repoPath, outputFormat, nil
}

func scanRepo(repoPath string, outputFormat report.Format) {
	if _, err := os.Stat(repoPath); err == nil {
		scanLocalRepo(repoPath, repoPath, outputFormat)
		return
	} else if !os.IsNotExist(err) {
		fmt.Printf("Error: cannot access path '%s': %v\n", repoPath, err)
		os.Exit(1)
	}

	repo, ok := parseGitHubRepo(repoPath)
	if !ok {
		fmt.Printf("Error: path '%s' does not exist\n", repoPath)
		os.Exit(1)
	}

	tempDir, err := cloneGitHubRepo(repo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	scanLocalRepo(tempDir, repo.DisplayName(), outputFormat)
}

func scanLocalRepo(repoPath, displayPath string, outputFormat report.Format) {
	// Check if path exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		fmt.Printf("Error: path '%s' does not exist\n", repoPath)
		os.Exit(1)
	}

	// Check if it's a git repo
	gitPath := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitPath); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Warning: .git directory not found. This may not be a Git repository.")
	}

	// Scan
	results := scan.Scan(repoPath)

	// Calculate score
	scoreAnalysis := score.Analyze(results)

	// Report
	output, err := report.Render(results, scoreAnalysis, displayPath, outputFormat)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(output)
}

type githubRepo struct {
	Owner string
	Name  string
}

func (repo githubRepo) CloneURL() string {
	return fmt.Sprintf("https://github.com/%s/%s.git", repo.Owner, repo.Name)
}

func (repo githubRepo) DisplayName() string {
	return fmt.Sprintf("github.com/%s/%s", repo.Owner, repo.Name)
}

func parseGitHubRepo(input string) (githubRepo, bool) {
	input = strings.TrimSpace(input)
	input = strings.TrimSuffix(input, "/")
	input = strings.TrimPrefix(input, "https://")
	input = strings.TrimPrefix(input, "http://")
	input = strings.TrimPrefix(input, "www.")
	input = strings.TrimPrefix(input, "github.com/")

	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return githubRepo{}, false
	}

	owner := strings.TrimSpace(parts[0])
	name := strings.TrimSuffix(strings.TrimSpace(parts[1]), ".git")
	if !isGitHubPathSegment(owner) || !isGitHubPathSegment(name) {
		return githubRepo{}, false
	}

	return githubRepo{Owner: owner, Name: name}, true
}

func isGitHubPathSegment(segment string) bool {
	if segment == "" || strings.HasPrefix(segment, ".") || strings.HasSuffix(segment, ".") {
		return false
	}

	for _, r := range segment {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= 'A' && r <= 'Z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '-' || r == '_' || r == '.' {
			continue
		}
		return false
	}

	return true
}

func cloneGitHubRepo(repo githubRepo) (string, error) {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return "", fmt.Errorf("git is not installed or not available in PATH")
	}

	tempDir, err := os.MkdirTemp("", "repoxray-github-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	cmd := exec.Command(gitPath, "clone", "--depth", "1", repo.CloneURL(), tempDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = err.Error()
		}
		return "", fmt.Errorf("failed to clone %s: %s", repo.DisplayName(), message)
	}

	return tempDir, nil
}

func printHelp() {
	fmt.Println("RepoXray - Repository Health Analyzer")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  repoxray scan <path|owner/name|github.com/owner/name> [--format text|json|markdown]")
	fmt.Println("  repoxray version                                   Print version information")
	fmt.Println("  repoxray help                                      Print this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  repoxray scan .")
	fmt.Println("  repoxray scan . --format json")
	fmt.Println("  repoxray scan . --format markdown")
	fmt.Println("  repoxray scan /path/to/repo")
	fmt.Println("  repoxray scan biomejs/biome")
	fmt.Println("  repoxray scan github.com/biomejs/biome")
}
