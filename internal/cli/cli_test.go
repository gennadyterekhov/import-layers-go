package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gennadyterekhov/import-layers-go/internal/analyzer"
)

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Setenv("GOWORK", "off")
	write(t, root, "go.mod", "module example.com/project\n\ngo 1.23.0\n")
	write(t, root, "high/high.go", "package high\n")
	write(t, root, "low/low.go", "package low\n")
	return root
}

func write(t *testing.T, root, name, contents string) {
	t.Helper()
	filename := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(filename), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runCLI(cwd string, args ...string) (int, string, string) {
	var stdout, stderr bytes.Buffer
	status := Run(cwd, args, &stdout, &stderr)
	return status, stdout.String(), stderr.String()
}

func TestDefaultConfigIsFoundAboveInvocationDirectory(t *testing.T) {
	root := fixture(t)
	write(t, root, "import_layers.json", `{"layers":["high","low"]}`)
	write(t, root, "import_layers.yaml", "not valid: [")
	status, stdout, stderr := runCLI(filepath.Join(root, "high"))
	if status != 0 || stdout != "OK\n" || stderr != "" {
		t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}
}

func TestViolationAndGraphUseProjectRelativePaths(t *testing.T) {
	root := fixture(t)
	write(t, root, "config/layers.yaml", "layers: [high, low]\ndebug: true\n")
	write(t, root, "high/high.go", "package high\n\nimport _ \"example.com/project/low\"\n")

	status, stdout, stderr := runCLI(root, "high", "-c", "config/layers.yaml", "--graph", "low")
	if status != 1 || !strings.Contains(stderr, "high/high.go:3 -> example.com/project/low: cannot import package from lower layer (high -> low)") ||
		!strings.Contains(stderr, "processing file high/high.go") {
		t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}
	var actual analyzer.Graph
	if err := json.Unmarshal([]byte(stdout), &actual); err != nil {
		t.Fatal(err)
	}
	if len(actual.Nodes) != 2 || actual.Nodes[0] != "high/high.go" || actual.Nodes[1] != "low/low.go" ||
		len(actual.Edges) != 1 || actual.Edges[0] != (analyzer.Edge{From: "high/high.go", To: "low/low.go"}) {
		t.Fatalf("unexpected graph: %+v", actual)
	}
}

func TestConfiguredRootAndLegacyIgnoreTests(t *testing.T) {
	root := fixture(t)
	write(t, root, "import_layers.yaml", "layers: [high, low]\nignore_tests: true\nrootDir: high\ndebug: true\n")
	write(t, root, "high/high_test.go", "package high\n\nimport _ \"example.com/project/low\"\n")
	status, stdout, stderr := runCLI(root, "low")
	if status != 0 || stdout != "processing file high/high.go\nOK\n" || stderr != "" {
		t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}
}

func TestUnresolvedImportFailsAnalysis(t *testing.T) {
	root := fixture(t)
	write(t, root, "import_layers.json", `{"layers":["high","low"]}`)
	write(t, root, "high/high.go", "package high\n\nimport _ \"example.com/project/missing\"\n")
	status, stdout, stderr := runCLI(root)
	if status != 2 || stdout != "" || !strings.Contains(stderr, "example.com/project/missing") {
		t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout, stderr)
	}
}
