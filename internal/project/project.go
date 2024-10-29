package project

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectRoot() (string, error) {
	// Start from the current working directory.
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("unable to get current working directory: %w", err)
	}

	// Iterate up the directory tree until we find go.mod
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil { // go.mod exists
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)

		// Reached the root directory without finding go.mod
		if parentDir == currentDir {
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}

		currentDir = parentDir
	}
}
