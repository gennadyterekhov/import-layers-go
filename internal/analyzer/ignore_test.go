package analyzer

import (
	"strings"
	"testing"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/layers"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestIgnoreTestsKeepsCheckingProductionFiles(t *testing.T) {
	cfg := &config.Config{
		IgnoreTests: true,
		Layers:      layers.New([]string{"highignored", "highlogicandtest", "low"}),
	}
	results := analysistest.Run(t, analysistest.TestData(), New(cfg).Analyzer,
		"low", "highignored", "highlogicandtest")

	var count int
	for _, result := range results {
		for _, diagnostic := range result.Diagnostics {
			filename := result.Pass.Fset.Position(diagnostic.Pos).Filename
			if strings.HasSuffix(filename, "_test.go") {
				t.Errorf("unexpected test diagnostic in %s", filename)
			}
			count++
		}
	}
	if count == 0 {
		t.Fatal("no diagnostics for the production file")
	}
}
