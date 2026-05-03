package checks

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"repoxray/internal/types"
)

// DefaultChecks returns the standard repository diagnostics.
func DefaultChecks() []types.Check {
	return []types.Check{
		CheckReadme{},
		CheckLicense{},
		CheckContributing{},
		CheckCodeOfConduct{},
		CheckSecurity{},
		CheckCI{},
		CheckWorkflowPermissions{},
		CheckWorkflowWriteAllPermissions{},
		CheckWorkflowActionPins{},
		CheckWorkflowPullRequestTarget{},
		CheckWorkflowPullRequestTargetSecrets{},
		CheckIssueTemplates{},
		CheckPRTemplate{},
		CheckPackageFiles{},
		CheckTests{},
		CheckGit{},
	}
}

// CheckReadme checks if README.md exists
type CheckReadme struct{}

func (c CheckReadme) Run(ctx types.Context) types.Result {
	if ctx.HasPath("README.md") {
		return types.Result{
			ID:             "readme",
			Title:          "README.md",
			Description:    "A README.md file provides essential information about the project.",
			Status:         types.Pass,
			Message:        "README.md exists",
			Points:         10,
			MaxPoints:      10,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "readme",
		Title:          "README.md",
		Description:    "A README.md file provides essential information about the project.",
		Status:         types.Fail,
		Message:        "README.md is missing",
		Points:         0,
		MaxPoints:      10,
		Recommendation: "Add a README.md file to describe your project.",
	}
}

// CheckLicense checks if LICENSE exists
type CheckLicense struct{}

func (c CheckLicense) Run(ctx types.Context) types.Result {
	if licenseFile, ok := findLicenseFile(ctx); ok {
		return types.Result{
			ID:             "license",
			Title:          "LICENSE",
			Description:    "A LICENSE file defines the terms under which the project can be used.",
			Status:         types.Pass,
			Message:        fmt.Sprintf("License file exists: %s", licenseFile),
			Points:         10,
			MaxPoints:      10,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "license",
		Title:          "LICENSE",
		Description:    "A LICENSE file defines the terms under which the project can be used.",
		Status:         types.Fail,
		Message:        "LICENSE file is missing",
		Points:         0,
		MaxPoints:      10,
		Recommendation: "Add a LICENSE file, or a recognized license variant such as LICENSE-MIT or LICENSE-APACHE, so users know how they can use, modify, and distribute the project.",
	}
}

// IsLicenseFile reports whether name is a recognized root license filename.
func IsLicenseFile(name string) bool {
	normalized := strings.ToUpper(name)

	exact := map[string]bool{
		"LICENSE":     true,
		"LICENSE.MD":  true,
		"LICENSE.TXT": true,
		"LICENCE":     true,
		"LICENCE.MD":  true,
		"LICENCE.TXT": true,
		"COPYING":     true,
		"COPYING.MD":  true,
		"COPYING.TXT": true,
		"COPYRIGHT":   true,
		"UNLICENSE":   true,
	}

	if exact[normalized] {
		return true
	}

	return strings.HasPrefix(normalized, "LICENSE-") ||
		strings.HasPrefix(normalized, "LICENSE.") ||
		strings.HasPrefix(normalized, "LICENCE-") ||
		strings.HasPrefix(normalized, "LICENCE.") ||
		strings.HasPrefix(normalized, "COPYING-") ||
		strings.HasPrefix(normalized, "COPYING.")
}

func findLicenseFile(ctx types.Context) (string, bool) {
	rootFiles := rootFileNames(ctx)
	sort.Strings(rootFiles)

	for _, name := range rootFiles {
		if IsLicenseFile(name) {
			return name, true
		}
	}

	return "", false
}

func rootFileNames(ctx types.Context) []string {
	var names []string

	if ctx.Files != nil {
		for _, file := range ctx.Files {
			if file.IsDir || strings.Contains(file.Path, "/") {
				continue
			}
			names = append(names, file.Path)
		}
		return names
	}

	entries, err := os.ReadDir(ctx.RepoPath)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	return names
}

// CheckContributing checks if CONTRIBUTING.md exists
type CheckContributing struct{}

func (c CheckContributing) Run(ctx types.Context) types.Result {
	if ctx.HasPath("CONTRIBUTING.md") {
		return types.Result{
			ID:             "contributing",
			Title:          "CONTRIBUTING.md",
			Description:    "A CONTRIBUTING.md file guides contributors on how to contribute.",
			Status:         types.Pass,
			Message:        "CONTRIBUTING.md exists",
			Points:         5,
			MaxPoints:      5,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "contributing",
		Title:          "CONTRIBUTING.md",
		Description:    "A CONTRIBUTING.md file guides contributors on how to contribute.",
		Status:         types.Warn,
		Message:        "CONTRIBUTING.md is missing",
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Consider adding a CONTRIBUTING.md file to help contributors.",
	}
}

// CheckCodeOfConduct checks if CODE_OF_CONDUCT.md exists
type CheckCodeOfConduct struct{}

func (c CheckCodeOfConduct) Run(ctx types.Context) types.Result {
	if ctx.HasPath("CODE_OF_CONDUCT.md") {
		return types.Result{
			ID:             "code_of_conduct",
			Title:          "CODE_OF_CONDUCT.md",
			Description:    "A CODE_OF_CONDUCT.md file sets expectations for behavior in the community.",
			Status:         types.Pass,
			Message:        "CODE_OF_CONDUCT.md exists",
			Points:         5,
			MaxPoints:      5,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "code_of_conduct",
		Title:          "CODE_OF_CONDUCT.md",
		Description:    "A CODE_OF_CONDUCT.md file sets expectations for behavior in the community.",
		Status:         types.Warn,
		Message:        "CODE_OF_CONDUCT.md is missing",
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Consider adding a CODE_OF_CONDUCT.md file.",
	}
}

// CheckSecurity checks if SECURITY.md exists
type CheckSecurity struct{}

func (c CheckSecurity) Run(ctx types.Context) types.Result {
	if ctx.HasPath("SECURITY.md") {
		return types.Result{
			ID:             "security",
			Title:          "SECURITY.md",
			Description:    "A SECURITY.md file provides security reporting guidelines.",
			Status:         types.Pass,
			Message:        "SECURITY.md exists",
			Points:         5,
			MaxPoints:      5,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "security",
		Title:          "SECURITY.md",
		Description:    "A SECURITY.md file provides security reporting guidelines.",
		Status:         types.Warn,
		Message:        "SECURITY.md is missing",
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Consider adding a SECURITY.md file for security disclosures.",
	}
}

// CheckCI checks if .github/workflows exists
type CheckCI struct{}

func (c CheckCI) Run(ctx types.Context) types.Result {
	workflowFiles := workflowFilePaths(ctx)
	if len(workflowFiles) > 0 {
		return types.Result{
			ID:             "ci",
			Title:          "CI Workflows",
			Description:    "Continuous Integration ensures code quality through automated testing.",
			Status:         types.Pass,
			Message:        fmt.Sprintf("CI workflows exist: %d workflow file(s)", len(workflowFiles)),
			Points:         15,
			MaxPoints:      15,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "ci",
		Title:          "CI Workflows",
		Description:    "Continuous Integration ensures code quality through automated testing.",
		Status:         types.Fail,
		Message:        "CI workflow files are missing",
		Points:         0,
		MaxPoints:      15,
		Recommendation: "Set up CI workflows in .github/workflows using .yml or .yaml files.",
	}
}

// CheckIssueTemplates checks if issue templates exist
type CheckIssueTemplates struct{}

func (c CheckIssueTemplates) Run(ctx types.Context) types.Result {
	if ctx.HasPath(".github/ISSUE_TEMPLATE") {
		return types.Result{
			ID:             "issue_templates",
			Title:          "Issue Templates",
			Description:    "Issue templates guide users in reporting issues effectively.",
			Status:         types.Pass,
			Message:        "Issue templates exist",
			Points:         5,
			MaxPoints:      5,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "issue_templates",
		Title:          "Issue Templates",
		Description:    "Issue templates guide users in reporting issues effectively.",
		Status:         types.Warn,
		Message:        "Issue templates are missing",
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Consider adding issue templates in .github/ISSUE_TEMPLATE/.",
	}
}

// CheckPRTemplate checks if pull request template exists
type CheckPRTemplate struct{}

func (c CheckPRTemplate) Run(ctx types.Context) types.Result {
	if ctx.HasPath(".github/PULL_REQUEST_TEMPLATE.md") {
		return types.Result{
			ID:             "pr_template",
			Title:          "Pull Request Template",
			Description:    "A PR template ensures consistent and complete pull requests.",
			Status:         types.Pass,
			Message:        "Pull request template exists",
			Points:         5,
			MaxPoints:      5,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "pr_template",
		Title:          "Pull Request Template",
		Description:    "A PR template ensures consistent and complete pull requests.",
		Status:         types.Warn,
		Message:        "Pull request template is missing",
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Consider adding a pull request template in .github/PULL_REQUEST_TEMPLATE.md.",
	}
}

// CheckPackageFiles checks for common package files
type CheckPackageFiles struct{}

func (c CheckPackageFiles) Run(ctx types.Context) types.Result {
	packageFiles := []string{"go.mod", "package.json", "Cargo.toml", "pyproject.toml"}
	found := false
	for _, file := range packageFiles {
		if ctx.HasPath(file) {
			found = true
			break
		}
	}
	if found {
		return types.Result{
			ID:             "package_files",
			Title:          "Package Files",
			Description:    "Package files indicate the project is properly configured for its language.",
			Status:         types.Pass,
			Message:        "Package files exist",
			Points:         10,
			MaxPoints:      10,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "package_files",
		Title:          "Package Files",
		Description:    "Package files indicate the project is properly configured for its language.",
		Status:         types.Fail,
		Message:        "Package files are missing",
		Points:         0,
		MaxPoints:      10,
		Recommendation: "Add appropriate package files (e.g., go.mod for Go projects).",
	}
}

// CheckTests checks if the repository appears to have tests
type CheckTests struct{}

func (c CheckTests) Run(ctx types.Context) types.Result {
	if ctx.HasFileWithSuffix("_test.go") {
		return types.Result{
			ID:             "tests",
			Title:          "Tests",
			Description:    "Tests ensure code reliability and prevent regressions.",
			Status:         types.Pass,
			Message:        "Tests exist",
			Points:         15,
			MaxPoints:      15,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "tests",
		Title:          "Tests",
		Description:    "Tests ensure code reliability and prevent regressions.",
		Status:         types.Fail,
		Message:        "Tests are missing",
		Points:         0,
		MaxPoints:      15,
		Recommendation: "Add unit tests to verify your code works correctly.",
	}
}

// CheckGit checks if .git directory exists
type CheckGit struct{}

func (c CheckGit) Run(ctx types.Context) types.Result {
	if ctx.HasPath(".git") {
		return types.Result{
			ID:             "git",
			Title:          ".git Directory",
			Description:    "A .git directory indicates this is a Git repository.",
			Status:         types.Pass,
			Message:        ".git directory exists",
			Points:         5,
			MaxPoints:      5,
			Recommendation: "",
		}
	}
	return types.Result{
		ID:             "git",
		Title:          ".git Directory",
		Description:    "A .git directory indicates this is a Git repository.",
		Status:         types.Warn,
		Message:        ".git directory is missing",
		Points:         0,
		MaxPoints:      5,
		Recommendation: "Initialize a Git repository with 'git init'.",
	}
}
