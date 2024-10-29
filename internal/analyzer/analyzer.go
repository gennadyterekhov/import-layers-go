package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"golang.org/x/tools/go/analysis"
)

type Analyzer struct {
	Analyzer *analysis.Analyzer
	config   *config.Config
}

func New(config *config.Config) *Analyzer {
	inst := &Analyzer{
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

func (f *Analyzer) run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg == nil {
		return nil, nil
	}
	pkgPath := pass.Pkg.Path()

	currentLayer := f.config.GetLayer(pkgPath)

	if f.config.Debug() {
		fmt.Println("Analyzer.run", "pkgPath", pkgPath, "currentLayer", currentLayer)
	}

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

func (f *Analyzer) inspectImport(importNode *ast.ImportSpec, currentLayer int) (bool, token.Pos) {
	if importNode != nil && importNode.Path != nil {
		importedPkgPath := strings.Trim(importNode.Path.Value, "\"")

		layer := f.config.GetLayer(importedPkgPath)

		if f.config.Debug() && importNode != nil && importNode.Path != nil {
			fmt.Println("    importedPkgPath", importedPkgPath, "layer", layer)
		}

		if layer != 0 && layer < currentLayer {
			return false, importNode.Pos()
		}
	}

	return true, 0
}
