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
	err := runPipeline("--project", struct {
		Language   string
		Layers     []string
		Trees      map[string]json.RawMessage
		References []analyzer.Reference
	}{"go", result.Layers, result.Trees, result.References}, &checked)
	return checked, err
}

func parseProject(root string, files []string) (analyzer.Snapshot, error) {
	var snapshot analyzer.Snapshot
	err := runPipeline("--parse-project", struct {
		Language string
		Root     string
		Files    []string
	}{"go", root, files}, &snapshot)
	return snapshot, err
}

func runPipeline(mode string, request, result any) error {
	var root string
	if _, source, _, ok := runtime.Caller(0); ok {
		candidate := filepath.Join(filepath.Dir(source), "..", "..", "..")
		if _, err := os.Stat(filepath.Join(candidate, "degra", "go.mod")); err == nil {
			root = candidate
		}
	}
	var cmd *exec.Cmd
	if directory := os.Getenv("IMPORT_LAYERS_TOOLS"); directory != "" {
		cmd = exec.Command(filepath.Join(directory, "degra"), mode, root)
	} else if binary, err := exec.LookPath("degra"); err == nil {
		cmd = exec.Command(binary, mode, root)
	} else if root != "" {
		cmd = exec.Command("go", "-C", filepath.Join(root, "degra"), "run", ".", mode, root)
	} else {
		return fmt.Errorf("degra not found on PATH")
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	cmd.Stdin = bytes.NewReader(payload)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	data, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("degra: %s%w", stderr.String(), err)
	}
	return json.Unmarshal(data, result)
}
