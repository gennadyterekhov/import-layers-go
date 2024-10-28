package collector

import (
	"testing"

	"github.com/gennadyterekhov/levelslib/internal/data"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestCanCollectData(t *testing.T) {
	commonData := data.New()
	analyzer := New(commonData)
	analysistest.Run(t, analysistest.TestData(), analyzer.Analyzer, "low", "high")
	assert.Equal(t, "low", analyzer.commonData.Deptree()["high"][0])
}
