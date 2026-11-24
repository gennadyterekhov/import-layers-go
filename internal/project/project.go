package project

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func FindRoot(start string) (string, error) {
	for dir := start; ; dir = filepath.Dir(dir) {
		if info, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && info.Mode().IsRegular() {
			return dir, nil
		}
		if filepath.Dir(dir) == dir {
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}
	}
}

type Root string

func (root Root) RelativePath(filename string) (string, bool) {
	relative, err := filepath.Rel(string(root), filename)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", false
	}
	relative = filepath.ToSlash(relative)
	for dir := path.Dir(relative); dir != "."; dir = path.Dir(dir) {
		switch path.Base(dir) {
		case ".git", "vendor":
			return "", false
		}
	}
	return relative, true
}

func IsTestFile(filename string) bool {
	return strings.HasSuffix(filepath.Base(filename), "_test.go")
}
