package filechecker

import (
	"go/ast"
	"go/token"
	"strconv"
)

type Import struct {
	Path     string
	Position token.Pos
}

func Imports(file *ast.File) ([]Import, []error) {
	imports := make([]Import, 0, len(file.Imports))
	var errs []error
	for _, spec := range file.Imports {
		imported, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		imports = append(imports, Import{Path: imported, Position: spec.Pos()})
	}
	return imports, errs
}
