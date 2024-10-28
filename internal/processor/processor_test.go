package processor

import (
	"testing"

	"github.com/gennadyterekhov/levelslib/internal/data"
	"github.com/stretchr/testify/assert"
)

func TestCanDetectInvalidImport(t *testing.T) {
	commonData := data.New()
	commonData.AddPackage("domain", 10)
	commonData.AddPackage("usecase", 5)
	commonData.AddPackage("infra", 1)

	commonData.AddImport("infra", "usecase")
	commonData.AddImport("usecase", "domain")
	commonData.AddImport("domain", "usecase") // invalid

	err := ProcessImports(commonData)
	assert.Error(t, err)
}

func TestDoesNotDetectInvalidImportIfOk(t *testing.T) {
	commonData := data.New()
	commonData.AddPackage("domain", 10)
	commonData.AddPackage("usecase", 5)
	commonData.AddPackage("infra", 1)

	commonData.AddImport("infra", "usecase")
	commonData.AddImport("usecase", "domain")

	err := ProcessImports(commonData)
	assert.NoError(t, err)
}
