package project

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller information")
	}

	// Iterate up the directory tree until we find go.mod
	currentDir := filepath.Dir(filename)
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
