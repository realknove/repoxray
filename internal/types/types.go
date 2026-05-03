package types

import (
	"os"
	"path/filepath"
	"strings"
)

type Status string

const (
	Pass Status = "pass"
	Warn Status = "warn"
	Fail Status = "fail"
)

type File struct {
	Path  string
	IsDir bool
}

type Context struct {
	RepoPath string
	Files    map[string]File
}

func (ctx Context) HasPath(path string) bool {
	path = filepath.ToSlash(filepath.Clean(path))

	if ctx.Files != nil {
		_, ok := ctx.Files[path]
		return ok
	}

	_, err := os.Stat(filepath.Join(ctx.RepoPath, path))
	return err == nil
}

func (ctx Context) HasFileWithSuffix(suffix string) bool {
	if ctx.Files != nil {
		for _, file := range ctx.Files {
			if !file.IsDir && strings.HasSuffix(file.Path, suffix) {
				return true
			}
		}
		return false
	}

	found := false
	_ = filepath.WalkDir(ctx.RepoPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return nil
		}

		if strings.HasSuffix(entry.Name(), suffix) {
			found = true
			return filepath.SkipAll
		}

		return nil
	})

	return found
}

type Check interface {
	Run(ctx Context) Result
}

type Result struct {
	ID             string
	Title          string
	Description    string
	Status         Status
	Message        string
	Points         int
	MaxPoints      int
	Recommendation string
}
