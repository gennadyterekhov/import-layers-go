package processor

import (
	"fmt"

	"github.com/gennadyterekhov/levelslib/internal/data"
)

func ProcessImports(data *data.CommonData) error {
	tempLevel := 0
	ok := false
	currentLevel := 0
	for pkgName, deps := range data.Deptree() {
		tempLevel, ok = data.PakToLev()[pkgName]
		if ok {
			currentLevel = tempLevel
		}

		for _, fullDepPath := range deps {

			tempLevel, ok = data.PakToLev()[fullDepPath] // error - we only have fnames as keys

			if ok && tempLevel < currentLevel {

				return errorLowImportedFromHigh(pkgName, currentLevel, fullDepPath, tempLevel)
			}
		}
	}

	return nil
}

func errorLowImportedFromHigh(pkgName string, currentLevel int, depName string, tempLevel int) error {
	return fmt.Errorf("wrong level. cannot import low from high. %s (%d) depends on %s (%d)", pkgName, currentLevel, depName, tempLevel)
}
