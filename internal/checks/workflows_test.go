package checks

import (
	"os"
	"path/filepath"
	"testing"

	"repoxray/internal/types"
)

func TestCheckCIDetectsWorkflowFiles(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"ci.yml": "name: CI\non: push\n",
	})

	result := CheckCI{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Pass {
		t.Fatalf("Expected Pass, got %v", result.Status)
	}
	if result.Message != "CI workflows exist: 1 workflow file(s)" {
		t.Fatalf("Expected workflow count message, got %q", result.Message)
	}
}

func TestCheckCIFailsWithoutWorkflowFiles(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, ".github", "workflows"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, ".github", "workflows", "notes.txt"), []byte("not yaml"), 0644); err != nil {
		t.Fatal(err)
	}

	result := CheckCI{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Fail {
		t.Fatalf("Expected Fail, got %v", result.Status)
	}
}

func TestWorkflowPermissionsPassWhenExplicit(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"ci.yml": "name: CI\non: push\npermissions:\n  contents: read\n",
	})

	result := CheckWorkflowPermissions{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Pass {
		t.Fatalf("Expected Pass, got %v", result.Status)
	}
}

func TestWorkflowPermissionsWarnWhenMissing(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"ci.yml": "name: CI\non: push\njobs:\n  test:\n    runs-on: ubuntu-latest\n",
	})

	result := CheckWorkflowPermissions{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Warn {
		t.Fatalf("Expected Warn, got %v", result.Status)
	}
}

func TestWorkflowWriteAllPermissionsFails(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"release.yaml": "name: Release\non: push\npermissions: write-all\n",
	})

	result := CheckWorkflowWriteAllPermissions{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Fail {
		t.Fatalf("Expected Fail, got %v", result.Status)
	}
}

func TestWorkflowActionPinsPassForVersionTagOrSHA(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"ci.yml": `name: CI
on: push
permissions:
  contents: read
jobs:
  test:
    steps:
      - uses: actions/checkout@v4
      - uses: docker://alpine:3.20
      - uses: ./local-action
      - uses: actions/setup-go@0db9f9b8f9ef5f9378c0d70f029a1db8d2ab50d6
`,
	})

	result := CheckWorkflowActionPins{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Pass {
		t.Fatalf("Expected Pass, got %v: %s", result.Status, result.Message)
	}
}

func TestWorkflowActionPinsWarnForUnpinnedAction(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"ci.yml": `name: CI
on: push
permissions:
  contents: read
jobs:
  test:
    steps:
      - uses: actions/checkout
`,
	})

	result := CheckWorkflowActionPins{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Warn {
		t.Fatalf("Expected Warn, got %v", result.Status)
	}
}

func TestWorkflowPullRequestTargetWarns(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"pr.yml": "name: PR\non: pull_request_target\npermissions:\n  contents: read\n",
	})

	result := CheckWorkflowPullRequestTarget{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Warn {
		t.Fatalf("Expected Warn, got %v", result.Status)
	}
}

func TestWorkflowPullRequestTargetSecretsFails(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"pr.yml": `name: PR
on: pull_request_target
permissions:
  contents: read
jobs:
  test:
    env:
      TOKEN: ${{ secrets.GITHUB_TOKEN }}
`,
	})

	result := CheckWorkflowPullRequestTargetSecrets{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Fail {
		t.Fatalf("Expected Fail, got %v", result.Status)
	}
}

func TestWorkflowPullRequestTargetSecretsPassesWithoutSecrets(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"pr.yml": "name: PR\non: pull_request_target\npermissions:\n  contents: read\n",
	})

	result := CheckWorkflowPullRequestTargetSecrets{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Pass {
		t.Fatalf("Expected Pass, got %v", result.Status)
	}
}

func TestWorkflowFilesIgnoreNestedAndNonYamlFiles(t *testing.T) {
	tempDir := newWorkflowFixture(t, map[string]string{
		"ci.yml":    "name: CI\non: push\n",
		"notes.txt": "not a workflow",
	})

	nestedDir := filepath.Join(tempDir, ".github", "workflows", "nested")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "nested.yml"), []byte("name: nested"), 0644); err != nil {
		t.Fatal(err)
	}

	paths := workflowFilePaths(types.Context{RepoPath: tempDir})
	if len(paths) != 1 || paths[0] != ".github/workflows/ci.yml" {
		t.Fatalf("Expected only root workflow file, got %#v", paths)
	}
}

func newWorkflowFixture(t *testing.T, files map[string]string) string {
	t.Helper()

	tempDir := t.TempDir()
	workflowsDir := filepath.Join(tempDir, ".github", "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		t.Fatal(err)
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(workflowsDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	return tempDir
}
