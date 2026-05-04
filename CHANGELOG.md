# Changelog

All notable changes to RepoXray will be documented in this file.

## v0.1.0 - 2026-05-04

Initial release.

### Added

- Local repository scanning.
- Public GitHub repository scanning by `owner/name` or `github.com/owner/name`.
- Weighted repository health scoring across seven categories:
  - Documentation
  - Licensing
  - Contributor readiness
  - CI and automation
  - Security posture
  - Maintenance signals
  - Project structure
- Text, JSON, and Markdown report output.
- Checks for README, license, contribution docs, code of conduct, security policy, CI workflows, workflow security, issue templates, pull request template, package files, tests, and Git metadata.
- CLI version command.
- Installation via `go install github.com/realknove/repoxray/cmd/repoxray@latest`.
