package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/realknove/repoxray/internal/types"
)

// CheckWorkflowPermissions checks whether workflow files define permissions.
type CheckWorkflowPermissions struct{}

func (c CheckWorkflowPermissions) Run(ctx types.Context) types.Result {
	workflows := readWorkflows(ctx)
	if len(workflows) == 0 {
		return types.Result{
			ID:             "workflow_permissions",
			Title:          "Workflow Permissions",
			Description:    "Explicit workflow permissions reduce the default token privileges used by GitHub Actions.",
			Status:         types.Warn,
			Message:        "Workflow permissions could not be checked because no workflow files were found",
			Points:         0,
			MaxPoints:      5,
			Recommendation: "Add workflow files and set explicit permissions for GitHub Actions.",
		}
	}

	var missing []string
	for _, workflow := range workflows {
		if !workflow.hasExplicitPermissions() {
			missing = append(missing, workflow.name)
		}
	}

	if len(missing) == 0 {
		return types.Result{
			ID:             "workflow_permissions",
			Title:          "Workflow Permissions",
			Description:    "Explicit workflow permissions reduce the default token privileges used by GitHub Actions.",
			Status:         types.Pass,
			Message:        "Workflow permissions are explicitly configured",
			Points:         5,
			MaxPoints:      5,
			Recommendation: "",
		}
	}

	return types.Result{
		ID:             "workflow_permissions",
		Title:          "Workflow Permissions",
		Description:    "Explicit workflow permissions reduce the default token privileges used by GitHub Actions.",
		Status:         types.Warn,
		Message:        fmt.Sprintf("Workflow permissions are missing in: %s", strings.Join(missing, ", ")),
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Set explicit permissions in each GitHub Actions workflow, preferably the least privileges required.",
	}
}

// CheckWorkflowWriteAllPermissions checks for broad write-all permissions.
type CheckWorkflowWriteAllPermissions struct{}

func (c CheckWorkflowWriteAllPermissions) Run(ctx types.Context) types.Result {
	workflows := readWorkflows(ctx)
	if len(workflows) == 0 {
		return workflowSecurityPass("workflow_write_all_permissions", "Workflow write-all permissions are not used")
	}

	var offenders []string
	for _, workflow := range workflows {
		if workflow.usesWriteAllPermissions() {
			offenders = append(offenders, workflow.name)
		}
	}

	if len(offenders) == 0 {
		return workflowSecurityPass("workflow_write_all_permissions", "Workflow write-all permissions are not used")
	}

	return types.Result{
		ID:             "workflow_write_all_permissions",
		Title:          "Workflow write-all Permissions",
		Description:    "The write-all permission grants broad write access to the GitHub token.",
		Status:         types.Fail,
		Message:        fmt.Sprintf("Workflow uses permissions: write-all in: %s", strings.Join(offenders, ", ")),
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Replace permissions: write-all with the narrowest explicit permissions required by each workflow.",
	}
}

// CheckWorkflowActionPins checks whether actions referenced by uses are pinned.
type CheckWorkflowActionPins struct{}

func (c CheckWorkflowActionPins) Run(ctx types.Context) types.Result {
	workflows := readWorkflows(ctx)
	if len(workflows) == 0 {
		return workflowSecurityPass("workflow_action_pins", "No unpinned GitHub Actions were found")
	}

	var unpinned []string
	for _, workflow := range workflows {
		for _, action := range workflow.unpinnedActions() {
			unpinned = append(unpinned, fmt.Sprintf("%s (%s)", workflow.name, action))
		}
	}

	if len(unpinned) == 0 {
		return workflowSecurityPass("workflow_action_pins", "GitHub Actions references are pinned by version, tag, or SHA")
	}

	return types.Result{
		ID:             "workflow_action_pins",
		Title:          "Workflow Action Pins",
		Description:    "Actions referenced with uses should include a version, tag, or SHA.",
		Status:         types.Warn,
		Message:        fmt.Sprintf("Unpinned GitHub Actions found: %s", strings.Join(unpinned, ", ")),
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Pin each GitHub Action reference with @version, @tag, or @SHA instead of using an unqualified action name.",
	}
}

// CheckWorkflowPullRequestTarget checks whether pull_request_target is used.
type CheckWorkflowPullRequestTarget struct{}

func (c CheckWorkflowPullRequestTarget) Run(ctx types.Context) types.Result {
	workflows := readWorkflows(ctx)
	var offenders []string
	for _, workflow := range workflows {
		if workflow.usesPullRequestTarget() {
			offenders = append(offenders, workflow.name)
		}
	}

	if len(offenders) == 0 {
		return workflowSecurityPass("workflow_pull_request_target", "pull_request_target is not used in workflows")
	}

	return types.Result{
		ID:             "workflow_pull_request_target",
		Title:          "Workflow pull_request_target",
		Description:    "pull_request_target runs with elevated repository context and should be used carefully.",
		Status:         types.Warn,
		Message:        fmt.Sprintf("pull_request_target is used in: %s", strings.Join(offenders, ", ")),
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Review workflows using pull_request_target and prefer pull_request unless elevated repository context is required.",
	}
}

// CheckWorkflowPullRequestTargetSecrets checks for secrets in pull_request_target workflows.
type CheckWorkflowPullRequestTargetSecrets struct{}

func (c CheckWorkflowPullRequestTargetSecrets) Run(ctx types.Context) types.Result {
	workflows := readWorkflows(ctx)
	var offenders []string
	for _, workflow := range workflows {
		if workflow.usesPullRequestTarget() && workflow.referencesSecrets() {
			offenders = append(offenders, workflow.name)
		}
	}

	if len(offenders) == 0 {
		return workflowSecurityPass("workflow_pull_request_target_secrets", "Secrets are not referenced from pull_request_target workflows")
	}

	return types.Result{
		ID:             "workflow_pull_request_target_secrets",
		Title:          "Workflow pull_request_target Secrets",
		Description:    "Secrets in pull_request_target workflows can expose sensitive credentials when handling untrusted changes.",
		Status:         types.Fail,
		Message:        fmt.Sprintf("Secrets appear in pull_request_target workflows: %s", strings.Join(offenders, ", ")),
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Remove secrets from pull_request_target workflows or redesign the workflow so untrusted pull requests cannot access sensitive values.",
	}
}

type workflowFile struct {
	path    string
	name    string
	content string
}

func workflowSecurityPass(id, message string) types.Result {
	return types.Result{
		ID:             id,
		Title:          "Workflow Security",
		Description:    "GitHub Actions workflow security checks help avoid risky CI configuration.",
		Status:         types.Pass,
		Message:        message,
		Points:         5,
		MaxPoints:      5,
		Recommendation: "",
	}
}

func workflowFilePaths(ctx types.Context) []string {
	var paths []string

	if ctx.Files != nil {
		for _, file := range ctx.Files {
			if file.IsDir || !isWorkflowFilePath(file.Path) {
				continue
			}
			paths = append(paths, file.Path)
		}
		sort.Strings(paths)
		return paths
	}

	workflowsDir := filepath.Join(ctx.RepoPath, ".github", "workflows")
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		relPath := filepath.ToSlash(filepath.Join(".github", "workflows", entry.Name()))
		if isWorkflowFilePath(relPath) {
			paths = append(paths, relPath)
		}
	}

	sort.Strings(paths)
	return paths
}

func isWorkflowFilePath(path string) bool {
	path = filepath.ToSlash(path)
	if strings.Count(path, "/") != 2 || !strings.HasPrefix(path, ".github/workflows/") {
		return false
	}

	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".yml") || strings.HasSuffix(lower, ".yaml")
}

func readWorkflows(ctx types.Context) []workflowFile {
	paths := workflowFilePaths(ctx)
	workflows := make([]workflowFile, 0, len(paths))

	for _, path := range paths {
		content, err := os.ReadFile(filepath.Join(ctx.RepoPath, filepath.FromSlash(path)))
		if err != nil {
			continue
		}

		workflows = append(workflows, workflowFile{
			path:    path,
			name:    filepath.Base(path),
			content: string(content),
		})
	}

	return workflows
}

func (workflow workflowFile) hasExplicitPermissions() bool {
	for _, line := range workflow.significantLines() {
		if strings.HasPrefix(strings.ToLower(line), "permissions:") {
			return true
		}
	}
	return false
}

func (workflow workflowFile) usesWriteAllPermissions() bool {
	for _, line := range workflow.significantLines() {
		normalized := strings.ToLower(strings.ReplaceAll(line, " ", ""))
		if strings.HasPrefix(normalized, "permissions:write-all") {
			return true
		}
	}
	return false
}

func (workflow workflowFile) unpinnedActions() []string {
	var unpinned []string
	for _, line := range workflow.significantLines() {
		action, ok := workflowUsesValue(line)
		if !ok {
			continue
		}

		if action == "" || strings.HasPrefix(action, "./") || strings.HasPrefix(action, ".github/") || strings.HasPrefix(action, "docker://") {
			continue
		}
		if !strings.Contains(action, "@") {
			unpinned = append(unpinned, action)
		}
	}
	return unpinned
}

func workflowUsesValue(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	lower := strings.ToLower(trimmed)

	switch {
	case strings.HasPrefix(lower, "uses:"):
		return strings.Trim(strings.TrimSpace(trimmed[len("uses:"):]), `"'`), true
	case strings.HasPrefix(lower, "- uses:"):
		return strings.Trim(strings.TrimSpace(trimmed[len("- uses:"):]), `"'`), true
	default:
		return "", false
	}
}

func (workflow workflowFile) usesPullRequestTarget() bool {
	for _, line := range workflow.significantLines() {
		if strings.Contains(strings.ToLower(line), "pull_request_target") {
			return true
		}
	}
	return false
}

func (workflow workflowFile) referencesSecrets() bool {
	for _, line := range workflow.significantLines() {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "secrets.") || strings.Contains(lower, "secrets[") {
			return true
		}
	}
	return false
}

func (workflow workflowFile) significantLines() []string {
	var lines []string
	for _, line := range strings.Split(workflow.content, "\n") {
		line = stripWorkflowComment(line)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}

func stripWorkflowComment(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "#") {
		return ""
	}

	if idx := strings.Index(line, " #"); idx >= 0 {
		return line[:idx]
	}

	return line
}
