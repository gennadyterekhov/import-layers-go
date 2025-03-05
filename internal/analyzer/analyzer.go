package analyzer

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/filechecker"
	"github.com/gennadyterekhov/import-layers-go/internal/reporter"
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
		fmt.Printf("\nAnalyzer.run pkgName: %v  pkgPath: %v  layer: %v\nignoreTests: %v\n", pass.Pkg.Name(), pkgPath, currentLayer, f.config.IgnoreTests())

		//err := ast.Print(pass.Fset, pass.Pkg)
		//if err != nil {
		//	return nil, fmt.Errorf("could not print ast")
		//}
	}

	if currentLayer == 0 { //can import anything from layer 0
		return nil, nil
	}

	//pkgImports := getPkgImports(pass.Pkg)
	//// TODO optimize to look for testing in the first place, in one traversal. implement after benchmark
	//if slices.Contains(pkgImports, "testing") {
	//	if f.config.IgnoreTests() {
	//		if f.config.Debug() {
	//			fmt.Println("    this fileSet contains a test , ignoring")
	//		}
	//		return nil, nil
	//	}
	//
	//	if f.config.Debug() {
	//		fmt.Println("    this fileSet contains a test")
	//	}
	//}

	rep := reporter.New(pass.Reportf)
	fchecker := filechecker.New(rep, currentLayer, f.config)

	for _, file := range pass.Files {
		fchecker.CheckFile(file)
	}

	return nil, nil
}

func getPkgImports(file *types.Package) []string {
	empty := make([]string, 0)
	if file == nil {
		return empty
	}

	for _, v := range file.Imports() {
		if v == nil {
			continue
		}

		empty = append(empty, strings.Trim(v.Path(), "\""))
	}

	return empty
}
