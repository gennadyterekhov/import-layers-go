package project

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type PackageSelection struct {
	Patterns []string
	files    map[string]map[string]bool
}

func (selection PackageSelection) IncludesFile(relative string) bool {
	files, restricted := selection.files[path.Dir(relative)]
	return !restricted || files[relative]
}

func (root Root) SelectPackages(cwd, rootDir string, paths []string) (PackageSelection, error) {
	configured := rootDir != ""
	if configured {
		cwd = string(root)
		paths = []string{rootDir}
	}
	if len(paths) == 0 {
		return PackageSelection{Patterns: []string{"./..."}}, nil
	}

	selection := PackageSelection{files: make(map[string]map[string]bool)}
	var packages, recursiveDirs []string
	for _, target := range paths {
		recursive := strings.HasSuffix(target, "/...") || target == "..."
		filesystemPath := target
		if recursive {
			filesystemPath = strings.TrimSuffix(target, "...")
		}
		if !filepath.IsAbs(filesystemPath) {
			filesystemPath = filepath.Join(cwd, filesystemPath)
		}
		info, err := os.Stat(filesystemPath)
		if err != nil {
			if configured || strings.HasPrefix(target, ".") || filepath.IsAbs(target) || filepath.Ext(target) == ".go" {
				return selection, fmt.Errorf("could not read path %s: %w", target, err)
			}
			// Preserve native import path and wildcard patterns.
			selection.Patterns = append(selection.Patterns, target)
			packages = append(packages, target)
			continue
		}

		absolute, err := filepath.EvalSymlinks(filesystemPath)
		if err != nil {
			return selection, err
		}
		relative, ok := root.RelativePath(absolute)
		if !ok {
			return selection, fmt.Errorf("path is outside the project or excluded: %s", target)
		}

		switch {
		case !info.IsDir():
			if configured || path.Ext(relative) != ".go" {
				return selection, fmt.Errorf("expected a directory or Go source file: %s", target)
			}
			directory := path.Dir(relative)
			if selection.files[directory] == nil {
				selection.files[directory] = make(map[string]bool)
			}
			selection.files[directory][relative] = true
			selection.Patterns = append(selection.Patterns, packagePattern(directory))
		case configured || recursive || !strings.HasPrefix(target, "./"):
			recursiveDirs = append(recursiveDirs, relative)
			selection.Patterns = append(selection.Patterns, packagePattern(relative)+"/...")
		default:
			packages = append(packages, relative)
			selection.Patterns = append(selection.Patterns, packagePattern(relative))
		}
	}

	for _, directory := range packages {
		delete(selection.files, directory)
	}
	for directory := range selection.files {
		for _, parent := range recursiveDirs {
			if parent == "." || directory == parent || strings.HasPrefix(directory, parent+"/") {
				delete(selection.files, directory)
				break
			}
		}
	}
	return selection, nil
}

func packagePattern(directory string) string {
	if directory == "." {
		return directory
	}
	return "./" + directory
}
