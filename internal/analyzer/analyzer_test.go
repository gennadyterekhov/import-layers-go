package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLevels(t *testing.T) {
	reportsDownlevel(t)
	doesNotReportUplevel(t)
}

func reportsDownlevel(t *testing.T) {
	var result analysistest.Result
	results := analysistest.Run(t, analysistest.TestData(), New(), "./downlevel")

	assert.Equal(t, 1, len(results))
	result = *results[0]
	assert.NoError(t, result.Err)

	assert.Equal(t, 4, len(result.Diagnostics))
}

func doesNotReportUplevel(t *testing.T) {
	var result analysistest.Result
	results := analysistest.Run(t, analysistest.TestData(), New(), "./uplevel")

	assert.Equal(t, 1, len(results))
	result = *results[0]
	assert.NoError(t, result.Err)

	assert.Equal(t, 0, len(result.Diagnostics))
}
