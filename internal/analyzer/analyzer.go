package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"golang.org/x/tools/go/analysis"
)

type Finalizer struct {
	Analyzer *analysis.Analyzer
	config   *config.Config
}

func New(config *config.Config) *Finalizer {
	inst := &Finalizer{
		config: config,
	}
	analyzer := &analysis.Analyzer{
		Name: "import_layers",
		Doc:  "check that higher layer packages do not depend on lower layer packages (dependency rule from clean architecture)",
		Run:  inst.run,
	}
	inst.Analyzer = analyzer

	return inst
}

func (f *Finalizer) run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg == nil {
		return nil, nil
	}
	pkgPath := pass.Pkg.Path()

	currentLayer := f.config.GetLayer(pkgPath)
	for _, file := range pass.Files {

		if file.Name == nil {
			return nil, nil
		}

		ast.Inspect(file, func(node ast.Node) bool {
			typedNode, ok := node.(*ast.File)
			if !ok {
				return true
			}

			for _, importNode := range typedNode.Imports {
				ok, pos := f.inspectImport(importNode, currentLayer)

				if !ok {
					pass.Reportf(
						pos,
						"wrong layer. cannot import %s from \"%s\"",
						importNode.Path.Value,
						pkgPath,
					)
				}
			}

			return true
		})
	}

	return nil, nil
}

func (f *Finalizer) inspectImport(importNode *ast.ImportSpec, currentLayer int) (bool, token.Pos) {
	if importNode != nil && importNode.Path != nil {
		importedPkgPath := strings.Trim(importNode.Path.Value, "\"")

		layer := f.config.GetLayer(importedPkgPath)

		if layer != 0 && layer < currentLayer {
			return false, importNode.Pos()
		}
	}
	return true, 0
}
