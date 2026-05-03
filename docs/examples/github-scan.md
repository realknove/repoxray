# Example Report: Public GitHub Repository

Command:

```bash
repoxray scan github.com/biomejs/biome
```

RepoXray clones the public repository into a temporary directory, scans the
local clone, and removes the temporary directory when the scan finishes.

Example output:

```text
RepoXray Repository Health Report
=================================

Summary
-------
Total checks: 16
Repository: github.com/biomejs/biome
Passed: 13, Warned: 3, Failed: 0

Score
-----
Maturity Score: 90/100
Rating: Excellent

Strengths
---------
[PASS] README.md exists
[PASS] CI workflows exist: 4 workflow file(s)
[PASS] Tests exist

Warnings
--------
[WARN] Workflow permissions are missing in: ci.yml

Critical Issues
---------------
No critical issues found.
```
