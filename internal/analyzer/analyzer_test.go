package analyzer

import (
	"testing"

	"github.com/gennadyterekhov/import-layers-go/internal/config"
	"github.com/gennadyterekhov/import-layers-go/internal/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestReportWhenLowImportedByHigh(t *testing.T) {
	var result analysistest.Result

	commonData := data.New()
	commonData.AddPackage("high", 10)
	commonData.AddPackage("low", 5)

	cfg := config.FromMap(commonData.PackageToLayer())

	results := analysistest.Run(t, analysistest.TestData(), New(cfg).Analyzer, "low", "high")

	require.Equal(t, 2, len(results))

	result = *results[0]

	var diag analysis.Diagnostic
	if len(result.Diagnostics) == 0 {
		require.Equal(t, 1, len(results[1].Diagnostics))
		diag = results[1].Diagnostics[0]
	} else {
		require.Equal(t, 1, len(results[0].Diagnostics))
		diag = results[0].Diagnostics[0]
	}

	assert.Equal(t, `cannot import package from lower layer`, diag.Message)
}
