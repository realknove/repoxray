package checks

import (
	"os"
	"path/filepath"
	"repoxray/internal/types"
	"testing"
)

func TestCheckReadme(t *testing.T) {
	// Create temp dir
	tempDir, err := os.MkdirTemp("", "repoxray_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	ctx := types.Context{RepoPath: tempDir}

	// Test without README
	result := CheckReadme{}.Run(ctx)
	if result.Status != types.Fail {
		t.Errorf("Expected Fail, got %v", result.Status)
	}

	// Create README.md
	readmePath := filepath.Join(tempDir, "README.md")
	if err := os.WriteFile(readmePath, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	result = CheckReadme{}.Run(ctx)
	if result.Status != types.Pass {
		t.Errorf("Expected Pass, got %v", result.Status)
	}
}

func TestCheckLicense(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{name: "LICENSE", file: "LICENSE"},
		{name: "LICENSE.md", file: "LICENSE.md"},
		{name: "LICENSE-MIT", file: "LICENSE-MIT"},
		{name: "LICENSE-APACHE", file: "LICENSE-APACHE"},
		{name: "LICENSE-APACHE-2.0", file: "LICENSE-APACHE-2.0"},
		{name: "LICENCE", file: "LICENCE"},
		{name: "COPYING", file: "COPYING"},
		{name: "COPYING.LESSER", file: "COPYING.LESSER"},
		{name: "UNLICENSE", file: "UNLICENSE"},
		{name: "lowercase license", file: "license"},
		{name: "lowercase license mit", file: "license-mit"},
		{name: "mixed case copying lesser", file: "Copying.Lesser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "repoxray_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			licensePath := filepath.Join(tempDir, tt.file)
			if err := os.WriteFile(licensePath, []byte("test"), 0644); err != nil {
				t.Fatal(err)
			}

			result := CheckLicense{}.Run(types.Context{RepoPath: tempDir})
			if result.Status != types.Pass {
				t.Fatalf("Expected Pass, got %v", result.Status)
			}

			expectedMessage := "License file exists: " + tt.file
			if result.Message != expectedMessage {
				t.Fatalf("Expected message %q, got %q", expectedMessage, result.Message)
			}
		})
	}
}

func TestCheckLicenseMissing(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "repoxray_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	result := CheckLicense{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Fail {
		t.Errorf("Expected Fail, got %v", result.Status)
	}
	if result.Message != "LICENSE file is missing" {
		t.Errorf("Expected missing license message, got %q", result.Message)
	}
}

func TestCheckLicenseOnlyChecksRepositoryRoot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "repoxray_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	nestedDir := filepath.Join(tempDir, "vendor", "dependency")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "LICENSE"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	result := CheckLicense{}.Run(types.Context{RepoPath: tempDir})
	if result.Status != types.Fail {
		t.Errorf("Expected Fail, got %v", result.Status)
	}
}

func TestIsLicenseFile(t *testing.T) {
	valid := []string{
		"LICENSE",
		"LICENSE.md",
		"LICENSE.txt",
		"LICENCE",
		"LICENCE.md",
		"LICENCE.txt",
		"COPYING",
		"COPYING.md",
		"COPYING.txt",
		"COPYRIGHT",
		"UNLICENSE",
		"LICENSE-MIT",
		"LICENSE-APACHE",
		"LICENSE-APACHE-2.0",
		"LICENSE.BSD",
		"COPYING.LESSER",
		"LICENCE-MIT",
		"license",
		"copying.lesser",
	}

	for _, name := range valid {
		t.Run(name, func(t *testing.T) {
			if !IsLicenseFile(name) {
				t.Fatalf("Expected %q to be recognized as a license file", name)
			}
		})
	}

	invalid := []string{
		"LICENSES",
		"LICENSE_TEMPLATE",
		"THIRD_PARTY_LICENSES",
		"docs/LICENSE",
		"package.json",
	}

	for _, name := range invalid {
		t.Run(name, func(t *testing.T) {
			if IsLicenseFile(name) {
				t.Fatalf("Expected %q not to be recognized as a license file", name)
			}
		})
	}
}

func TestCheckGit(t *testing.T) {
	// Create temp dir
	tempDir, err := os.MkdirTemp("", "repoxray_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	ctx := types.Context{RepoPath: tempDir}

	// Test without .git
	result := CheckGit{}.Run(ctx)
	if result.Status != types.Warn {
		t.Errorf("Expected Warn, got %v", result.Status)
	}

	// Create .git
	gitPath := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitPath, 0755); err != nil {
		t.Fatal(err)
	}

	result = CheckGit{}.Run(ctx)
	if result.Status != types.Pass {
		t.Errorf("Expected Pass, got %v", result.Status)
	}
}
