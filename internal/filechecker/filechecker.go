package filechecker

import (
	"fmt"
	"go/ast"
	"slices"
	"strings"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/reporter"
)

type FileChecker struct {
	currentLayer int // abstraction layer of the file being checked
	rep          *reporter.Reporter
	conf         *config.Config
}

func New(rep *reporter.Reporter, currentLayer int, conf *config.Config) *FileChecker {
	return &FileChecker{
		currentLayer: currentLayer,
		rep:          rep,
		conf:         conf,
	}
}

func (fc *FileChecker) CheckFile(file *ast.File) {
	if file == nil {
		return
	}
	if file.Name == nil {
		return
	}

	if fc.conf.Debug() {
		fmt.Println("  checking file ", file.Name.String())
	}

	// TODO optimize to look for testing in the first place, in one traversal. implement after benchmark
	if slices.Contains(getImports(file), "testing") {
		if fc.conf.IgnoreTests() {
			if fc.conf.Debug() {
				fmt.Println("    this file is a test , ignoring")
			}
			return
		}

		if fc.conf.Debug() {
			fmt.Println("    this file is a test")
		}
	}

	ast.Inspect(file, fc.inspectAst)
}

// see ast.Inspect in "go/ast"
func (fc *FileChecker) inspectAst(node ast.Node) bool {
	if node == nil {
		return false
	}

	typedNode, ok := node.(*ast.File)
	if !ok {
		return true
	}

	if fc.conf.Debug() {
		fmt.Println("    current node type is file. looking through imports")
	}

	importOk := false

	for _, importNode := range typedNode.Imports {

		if importNode != nil && importNode.Path != nil {
			importedPkgPath := strings.Trim(importNode.Path.Value, "\"")

			if fc.conf.IgnoreTests() && importedPkgPath == "testing" {
				if fc.conf.Debug() {
					fmt.Println("    caught import testing when ignoring tests. skipping file")
				}
				return true
			}

			importOk = checkImport(fc.currentLayer, importedPkgPath, fc.conf.GetLayer)
			if !importOk {
				fc.rep.AddReport(
					typedNode,
					importNode,
				)
			}

			if fc.conf.Debug() {
				if importOk {
					fmt.Println("    import ", importedPkgPath, " ✅")
				} else {
					fmt.Println("    import ", importedPkgPath, " ❌")
				}
			}

		}
	}

	return true
}

// checkImport accepts int and string because currentLayer is known in advance, and we don't want to compute it every time
func checkImport(currentLayer int, importedPkgPath string, getLayer func(string) int) bool {
	layer := getLayer(importedPkgPath)

	return layer == 0 || layer >= currentLayer
}

func getImports(file *ast.File) []string {
	empty := make([]string, 0)
	if file == nil {
		return empty
	}

	for _, v := range file.Imports {
		if v == nil {
			continue
		}
		if v.Path != nil {
			empty = append(empty, strings.Trim(v.Path.Value, "\""))
		}
	}

	return empty
}
