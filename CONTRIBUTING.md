# Contributing To RepoXray

Thanks for helping improve RepoXray. The project aims to stay small, practical,
and easy to understand.

## Development Setup

You need Go installed. Then run:

```bash
go test ./...
```

Useful make targets:

```bash
make fmt
make lint
make test
make run
```

## Pull Requests

Before opening a pull request:

- Keep changes focused.
- Add or update tests for behavior changes.
- Run `make fmt`, `make lint`, and `make test`.
- Update documentation when user-facing behavior changes.

## Adding Checks

Repository diagnostics should implement the `types.Check` interface and be added
to the default check list in `internal/checks`.

Prefer simple, deterministic checks with clear messages and actionable
recommendations. If a heuristic is intentionally limited, document that in the
code or tests.
