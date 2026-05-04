# RepoXray

[![CI](https://github.com/yourusername/repoxray/actions/workflows/ci.yml/badge.svg)](https://github.com/yourusername/repoxray/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Project Status](https://img.shields.io/badge/status-early--stage-blue)](#roadmap)

RepoXray is an open-source Go CLI tool that analyzes the health, maturity, and maintainability of open-source repositories. It scans a local Git repository and produces a practical diagnostic report about repository quality, contributor readiness, maintenance signals, security posture, and project maturity.

RepoXray favors simple, explainable checks over opaque scoring. It is useful for quickly reviewing a repository before contributing, maintaining, adopting, or publishing it.

## Installation

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/repoxray.git
   cd repoxray
   ```

2. Build the binary:
   ```bash
   go build -o repoxray ./cmd/repoxray
   ```

3. Move to your PATH (optional):
   ```bash
   sudo mv repoxray /usr/local/bin/
   ```

## Usage

```bash
repoxray scan <path|owner/name|github.com/owner/name> [--format text|json|markdown]
```

Scan the current directory:
```bash
repoxray scan .
```

Text output is the default. You can also select a format explicitly:
```bash
repoxray scan . --format text
```

Generate machine-readable JSON:
```bash
repoxray scan . --format json
```

Generate a Markdown report for GitHub issues, pull requests, or README reports:
```bash
repoxray scan . --format markdown
```

RepoXray scores repositories with weighted categories: Documentation, Licensing,
Contributor readiness, CI and automation, Security posture, Maintenance signals,
and Project structure. Reports include both the overall score and a category
breakdown with status and recommendations.

Scan a specific repository:
```bash
repoxray scan /path/to/repository
```

Scan a public GitHub repository:
```bash
repoxray scan biomejs/biome
repoxray scan github.com/biomejs/biome
```

When scanning a GitHub repository, RepoXray clones it into a temporary directory,
scans the local clone, and deletes the temporary directory after the scan. This
requires `git` to be installed and available in your `PATH`.

Other commands:
```bash
repoxray version     # Print version
repoxray help        # Print help
```

## Development

Common development commands:

```bash
make fmt    # Format Go code
make lint   # Run go vet
make test   # Run all tests
make run    # Scan this repository
```

## Example Output

<img width="867" height="563" alt="image" src="https://github.com/user-attachments/assets/a242ee24-0df7-474e-b670-b66d262a0072" />

Text:

```
RepoXray Repository Health Report
=================================

Summary
-------
Total checks: 16
Repository: .
Passed: 16, Warned: 0, Failed: 0

Score
-----
Maturity Score: 100/100
Rating: Excellent

Category Breakdown
------------------
Documentation: 15/15 (100%) - excellent
Licensing: 10/10 (100%) - excellent
Contributor readiness: 15/15 (100%) - excellent
CI and automation: 20/20 (100%) - excellent
Security posture: 20/20 (100%) - excellent
Maintenance signals: 10/10 (100%) - excellent
Project structure: 10/10 (100%) - excellent

Strengths
---------
[PASS] README.md exists
[PASS] License file exists: LICENSE
[PASS] CONTRIBUTING.md exists
[PASS] CODE_OF_CONDUCT.md exists
[PASS] SECURITY.md exists
[PASS] CI workflows exist: 2 workflow file(s)
[PASS] Workflow permissions are explicitly configured
[PASS] Workflow write-all permissions are not used
[PASS] GitHub Actions references are pinned by version, tag, or SHA
[PASS] pull_request_target is not used in workflows
[PASS] Secrets are not referenced from pull_request_target workflows
[PASS] Issue templates exist
[PASS] Pull request template exists
[PASS] Package files exist
[PASS] Tests exist
[PASS] .git directory exists

Warnings
--------
No warnings.

Critical Issues
---------------
No critical issues found.

Recommendations
---------------
No recommendations. Repository looks healthy.
```

JSON:

```bash
repoxray scan . --format json
```

```json
{
  "repository": ".",
  "summary": {
    "total": 16,
    "passed": 16,
    "warned": 0,
    "failed": 0
  },
  "score": {
    "maturity": 100,
    "max": 100,
    "rating": "Excellent"
  },
  "categories": [
    {
      "id": "documentation",
      "name": "Documentation",
      "score": 15,
      "max_score": 15,
      "percentage": 100,
      "status": "excellent",
      "recommendations": []
    }
  ],
  "checks": []
}
```

Markdown:

```bash
repoxray scan . --format markdown
```

```markdown
# RepoXray Repository Health Report

## Summary

| Metric | Value |
| --- | --- |
| Repository | `.` |
| Total checks | 16 |
| Passed | 16 |
| Warned | 0 |
| Failed | 0 |

## Score

**Maturity Score:** 100/100

**Rating:** Excellent

## Category Breakdown

| Category | Score | Percentage | Status |
| --- | --- | --- | --- |
| Documentation | 15/15 | 100% | excellent |
| Licensing | 10/10 | 100% | excellent |
```

More examples are available in:

- [Local repository scan](docs/examples/local-scan.md)
- [Public GitHub repository scan](docs/examples/github-scan.md)

## Roadmap

- [ ] Add more detailed checks (e.g., code quality metrics)
- [x] Support for multiple output formats (text, JSON, Markdown)
- [ ] HTML report output
- [ ] Integration with GitHub API for additional metrics
- [ ] Web interface for repository analysis
- [ ] Plugin system for custom checks
- [ ] Support for non-Git repositories

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
