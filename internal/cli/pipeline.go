package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gennadyterekhov/import-layers-go/internal/analyzer"
)

type checkedProject struct {
	Edges      []analyzer.Edge
	Violations []string
}

func checkGraph(result analyzer.Result) (checkedProject, error) {
	var checked checkedProject
	var root string
	if _, source, _, ok := runtime.Caller(0); ok {
		candidate := filepath.Join(filepath.Dir(source), "..", "..", "..")
		if _, err := os.Stat(filepath.Join(candidate, "degra", "go.mod")); err == nil {
			root = candidate
		}
	}
	var cmd *exec.Cmd
	if directory := os.Getenv("IMPORT_LAYERS_TOOLS"); directory != "" {
		cmd = exec.Command(filepath.Join(directory, "degra"), "--project", root)
	} else if binary, err := exec.LookPath("degra"); err == nil {
		cmd = exec.Command(binary, "--project", root)
	} else if root != "" {
		cmd = exec.Command("go", "-C", filepath.Join(root, "degra"), "run", ".", "--project", root)
	} else {
		return checked, fmt.Errorf("degra not found on PATH")
	}
	payload, err := json.Marshal(struct {
		Language   string
		Root       string
		Layers     []string
		Files      []string
		References []analyzer.Reference
	}{"go", result.Root, result.Layers, result.Files, result.References})
	if err != nil {
		return checked, err
	}
	cmd.Stdin = bytes.NewReader(payload)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	data, err := cmd.Output()
	if err != nil {
		return checked, fmt.Errorf("degra: %s%w", stderr.String(), err)
	}
	err = json.Unmarshal(data, &checked)
	return checked, err
}
