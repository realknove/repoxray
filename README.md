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
repoxray scan <path|owner/name|github.com/owner/name>
```

Scan the current directory:
```bash
repoxray scan .
```

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
repoxray version    # Print version
repoxray help       # Print help
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


## Roadmap

- [ ] Add more detailed checks (e.g., code quality metrics)
- [ ] Support for multiple output formats (JSON, HTML)
- [ ] Integration with GitHub API for additional metrics
- [ ] Web interface for repository analysis
- [ ] Plugin system for custom checks
- [ ] Support for non-Git repositories

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
