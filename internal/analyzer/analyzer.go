package analyzer

import (
	"fmt"
	"os"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/filechecker"
	"github.com/gennadyterekhov/import-layers-go/internal/layers"
	"github.com/gennadyterekhov/import-layers-go/internal/project"
	"golang.org/x/tools/go/analysis"
)

type Analyzer struct {
	Analyzer *analysis.Analyzer
	config   *config.Config
}

func New(config *config.Config) *Analyzer {
	inst := &Analyzer{config: config}
	inst.Analyzer = &analysis.Analyzer{
		Name: "import_layers",
		Doc:  "check that higher layer packages do not depend on lower layer packages (dependency rule from clean architecture)",
		Run:  inst.run,
	}
	return inst
}

func (a *Analyzer) run(pass *analysis.Pass) (any, error) {
	if pass.Pkg == nil {
		return nil, nil
	}
	source := a.config.Layers.Match(pass.Pkg.Path(), "")
	for _, file := range pass.Files {
		filename := pass.Fset.PositionFor(file.Pos(), false).Filename
		if a.config.IgnoreTests && project.IsTestFile(filename) {
			continue
		}
		if a.config.Debug {
			if relative, ok := a.config.Root.RelativePath(filename); ok {
				filename = relative
			}
			fmt.Fprintln(os.Stderr, "processing file "+filename)
		}

		imports, _ := filechecker.Imports(file)
		for _, ref := range imports {
			if source.Forbids(a.config.Layers.Match(ref.Path, "")) {
				pass.Reportf(ref.Position, layers.ViolationMessage)
			}
		}
	}
	return nil, nil
}
