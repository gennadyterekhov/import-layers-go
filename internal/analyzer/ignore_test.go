package analyzer

import (
	"fmt"
	"testing"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestDoesNotReportTestFileIfConfigured(t *testing.T) {
	commonData := data.New()
	commonData.AddPackage("highignored", 10)
	commonData.AddPackage("low", 5)

	cfg := config.FromMap(commonData.PackageToLayer()).SetIgnoreTests(true)

	results := analysistest.Run(t, analysistest.TestData(), New(cfg).Analyzer, "low", "highignored")

	require.Equal(t, 4, len(results))

	count := 0

	for _, v := range results {
		for _, d := range v.Diagnostics {
			count++
			assert.Equal(t, `cannot import package from lower layer`, d.Message)
		}
	}
	require.Equal(t, 0, count)
}

func TestCanIgnoreWhenWholeRun(t *testing.T) {
	t.Skip()

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

	for i, c := range cases {

		commonData := data.New()
		commonData.AddPackage("koko", 10)
		commonData.AddPackage("kook", 10)
		commonData.AddPackage("okko", 10)
		commonData.AddPackage("okok", 10)
		commonData.AddPackage("low", 5)
		cfg := config.FromMap(commonData.PackageToLayer()).SetIgnoreTests(c.ignoreTests)

		patterns := make([]string, 0)
		patterns = append(patterns, "low")

		if c.logicOk {
			if c.testOk {
				patterns = append(patterns, "all/okok")
			} else {
				patterns = append(patterns, "all/okko")
			}
		} else {
			if c.testOk {
				patterns = append(patterns, "all/kook")
			} else {
				patterns = append(patterns, "all/koko")
			}
		}

		results := analysistest.Run(t, analysistest.TestData(), New(cfg).Analyzer, patterns...)

		diagCnt := 0
		for _, v := range results {
			for _, d := range v.Diagnostics {
				fmt.Println()
				fmt.Println(d.Message)
				fmt.Println()
				if d.Message != "" {
					diagCnt++
				}
			}

		}
		require.Equal(t, c.expectedNumberOrReports, diagCnt)
	}
}

// see pkgpass.TestWrongWithIgnoreTest
func TestReportNormalFileInOnePackageWithATestWhenTestsIgnored(t *testing.T) {
	t.Skip()
	commonData := data.New()
	commonData.AddPackage("highlogicandtest", 10)
	commonData.AddPackage("low", 5)

	cfg := config.FromMap(commonData.PackageToLayer()).SetIgnoreTests(true)

	results := analysistest.Run(t, analysistest.TestData(), New(cfg).Analyzer, "highlogicandtest")
	count := 0
	for _, v := range results {
		for _, d := range v.Diagnostics {
			count++
			assert.Equal(t, `cannot import package from lower layer`, d.Message)
		}
	}
	assert.Equal(t, 1, count)
}
