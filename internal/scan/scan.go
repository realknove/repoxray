package scan

import (
	"io/fs"
	"path/filepath"

	"github.com/realknove/repoxray/internal/checks"
	"github.com/realknove/repoxray/internal/types"
)

// Scan runs all checks on the repository at repoPath
func Scan(repoPath string) []types.Result {
	return ScanWithChecks(repoPath, checks.DefaultChecks())
}

// ScanWithChecks runs the provided checks on the repository at repoPath.
func ScanWithChecks(repoPath string, checks []types.Check) []types.Result {
	var results []types.Result
	ctx := NewContext(repoPath)

	for _, check := range checks {
		results = append(results, check.Run(ctx))
	}

	return results
}

// NewContext discovers repository files once and shares them with checks.
func NewContext(repoPath string) types.Context {
	return types.Context{
		RepoPath: repoPath,
		Files:    discoverFiles(repoPath),
	}
}

func discoverFiles(repoPath string) map[string]types.File {
	files := make(map[string]types.File)

	_ = filepath.WalkDir(repoPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		relPath, err := filepath.Rel(repoPath, path)
		if err != nil || relPath == "." {
			return nil
		}

		relPath = filepath.ToSlash(relPath)
		files[relPath] = types.File{
			Path:  relPath,
			IsDir: entry.IsDir(),
		}

		return nil
	})

	return files
}
