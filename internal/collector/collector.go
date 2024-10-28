package collector

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"github.com/gennadyterekhov/levelslib/internal/data"
	"golang.org/x/tools/go/analysis"
)

type Collector struct {
	Analyzer   *analysis.Analyzer
	commonData *data.CommonData
}

func New(commonData *data.CommonData) *Collector {
	inst := &Collector{
		commonData: commonData,
	}
	analyzer := &analysis.Analyzer{
		Name: "levels_collector",
		Doc:  "collect data for levels check",
		Run:  inst.run,
	}
	inst.Analyzer = analyzer

	return inst
}

func (c *Collector) run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg == nil {
		return nil, nil
	}
	pkgPath := pass.Pkg.Path()

	for _, file := range pass.Files {

		if file.Name == nil {
			return nil, nil
		}

		ast.Inspect(file, func(node ast.Node) bool {
			typedNode, ok := node.(*ast.File)
			if !ok {
				return true
			}

			cont, err := c.inspectFile(pkgPath, typedNode)
			if err != nil {
				pass.Reportf(1, err.Error())
			}
			return cont
		})
	}

	return nil, nil
}

func (c *Collector) inspectFile(pkgPath string, node *ast.File) (bool, error) {
	for _, importNode := range node.Imports {
		if importNode != nil && importNode.Path != nil {
			c.commonData.AddImport(pkgPath, strings.Trim(importNode.Path.Value, "\""))
		}
	}

	for _, dec := range node.Decls {
		typedDec, ok := dec.(*ast.GenDecl)
		if !ok {
			continue
		}

		if typedDec.Tok.String() == "const" {
			for _, spec := range typedDec.Specs {
				typedSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				if len(typedSpec.Values) == 1 && len(typedSpec.Names) == 1 && typedSpec.Names[0].String() == "LevelslibLevel" {
					levelAsStr := typedSpec.Values[0].(*ast.BasicLit).Value

					levelAsInt, err := strconv.Atoi(levelAsStr)
					if err != nil {
						return false, err
					}

					c.commonData.AddPackage(pkgPath, levelAsInt)
					return false, nil
				}
			}

		}
	}

	return false, fmt.Errorf("const not found")
}
