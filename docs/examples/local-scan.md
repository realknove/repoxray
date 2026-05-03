# Example Report: Local Repository

Command:

```bash
repoxray scan .
```

Example output:

```text
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
