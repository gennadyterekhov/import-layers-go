package filechecker

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/gennadyterekhov/import-layers-go/internal/reporter"
	"github.com/stretchr/testify/require"
)

const correctHighTestSrc = `
		package high
		
		import (
			"testing"
		)
		
		func TestFn(t *testing.T) {
		}
    `

const correctHighSrc = `
		package high

		func Fn() {
		}
    `

func TestCanIgnoreTestFiles(t *testing.T) {
	// truth table
	// 0 means failure (incorrect import)
	// ignoreTests logic test  Reports
	//  0			0		0		2
	//  0			0		1		1
	//  0			1		0		1
	//  0			1		1		0
	//  1			0		0		2
	//  1			0		1		1
	//  1			1		0		0
	//  1			1		1		0

	type testCase struct {
		ignoreTests             bool
		logicOk                 bool
		testOk                  bool
		expectedNumberOrReports int
	}
	cases := []testCase{
		{false, false, false, 2},
		{false, false, true, 1},
		{false, true, false, 1},
		{false, true, true, 0},

		{true, false, false, 1},
		{true, false, true, 1},
		{true, true, false, 0},
		{true, true, true, 0},
	}

	var nodeHighTest, nodeHigh, nodeLow *ast.File
	var err error
	for _, c := range cases {
		cfg := getConfig().SetIgnoreTests(c.ignoreTests)
		rep := reporter.NewMock()
		pass := New(rep, cfg.GetLayer("high"), cfg)

		fset := token.NewFileSet()

		nodeLow, err = parser.ParseFile(fset, "low.go", lowSrc, 0)
		require.NoError(t, err)

		if c.logicOk {
			nodeHigh, err = parser.ParseFile(fset, "high.go", correctHighSrc, 0)
			require.NoError(t, err)
		} else {
			nodeHigh, err = parser.ParseFile(fset, "high.go", highSrc, 0)
			require.NoError(t, err)
		}

		if c.testOk {
			nodeHighTest, err = parser.ParseFile(fset, "high_test.go", correctHighTestSrc, 0)
			require.NoError(t, err)
		} else {
			nodeHighTest, err = parser.ParseFile(fset, "high_test.go", highTestSrc, 0)
			require.NoError(t, err)
		}

		pass.CheckFile(nodeLow)
		pass.CheckFile(nodeHigh)
		pass.CheckFile(nodeHighTest)
		require.Equal(t, c.expectedNumberOrReports, len(rep.GetReports()))
	}
}
