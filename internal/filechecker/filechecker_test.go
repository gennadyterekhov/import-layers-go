package filechecker

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/data"
	"github.com/gennadyterekhov/import-layers-go/internal/reporter"
	"github.com/stretchr/testify/require"
)

const highTestSrc = `
		package high
		
		import (
			"low"
			"testing"
		)
		
		func TestFn(t *testing.T) {
			low.Fn()
		}
    `

const highSrc = `
		package high
		
		import (
			"low"
		)
		
		func Fn() {
			low.Fn()
		}
    `

const lowSrc = `
		package low
		
		import "fmt"
		
		func Fn() {
			fmt.Println("low")
		}
    `

func TestCorrect(t *testing.T) {

	cfg := getConfig()

	fset := token.NewFileSet()

	nodeLow, err := parser.ParseFile(fset, "low.go", lowSrc, 0)
	require.NoError(t, err)
	nodeHigh, err := parser.ParseFile(fset, "high.go", correctHighSrc, 0)
	require.NoError(t, err)
	nodeHighTest, err := parser.ParseFile(fset, "high_test.go", correctHighTestSrc, 0)
	require.NoError(t, err)

	rep := reporter.NewMock()
	pass := New(rep, cfg.GetLayer("high"), cfg)
	pass.CheckFile(nodeLow)
	pass.CheckFile(nodeHigh)
	pass.CheckFile(nodeHighTest)

	require.Equal(t, 0, len(rep.GetReports()))
}

func TestWrongWithTest(t *testing.T) {
	cfg := getConfig()

	fset := token.NewFileSet()

	nodeLow, err := parser.ParseFile(fset, "low.go", lowSrc, 0)
	require.NoError(t, err)
	nodeHigh, err := parser.ParseFile(fset, "high.go", highSrc, 0)
	require.NoError(t, err)
	nodeHighTest, err := parser.ParseFile(fset, "high_test.go", highTestSrc, 0)
	require.NoError(t, err)

	rep := reporter.NewMock()
	pass := New(rep, cfg.GetLayer("high"), cfg)
	pass.CheckFile(nodeLow)
	pass.CheckFile(nodeHigh)
	pass.CheckFile(nodeHighTest)

	// one report for high.go
	// one report for high_test.go
	require.Equal(t, 2, len(rep.GetReports()))
}

func getConfig() *config.Config {
	commonData := data.New()
	commonData.AddPackage("high", 10)
	commonData.AddPackage("low", 5)
	cfg := config.FromMap(commonData.PackageToLayer())
	return cfg
}
