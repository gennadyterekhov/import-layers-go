package filechecker

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestImportsPreservePathsAndPositions(t *testing.T) {
	source := "package high\n\nimport (\n  _ \"low/path\"\n  alias `other/path`\n)\n"
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "high.go", source, parser.SkipObjectResolution)
	if err != nil {
		t.Fatal(err)
	}

	imports, errs := Imports(file)
	if len(errs) != 0 || len(imports) != 2 {
		t.Fatalf("imports=%v errors=%v", imports, errs)
	}
	wantPaths := []string{"low/path", "other/path"}
	for index, ref := range imports {
		position := fset.Position(ref.Position)
		if ref.Path != wantPaths[index] || position.Line != 4+index || !strings.HasSuffix(position.Filename, "high.go") {
			t.Errorf("unexpected import %q at %s", ref.Path, position)
		}
	}
}
