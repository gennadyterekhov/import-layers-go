package finalizer

import (
	"github.com/gennadyterekhov/levelslib/internal/collector"
	"github.com/gennadyterekhov/levelslib/internal/data"
	"github.com/gennadyterekhov/levelslib/internal/processor"
	"golang.org/x/tools/go/analysis"
)

type Finalizer struct {
	Analyzer      *analysis.Analyzer
	commonData    *data.CommonData
	dataCollector *collector.Collector
}

func New(commonData *data.CommonData, dataCollector *collector.Collector) *Finalizer {
	inst := &Finalizer{
		commonData:    commonData,
		dataCollector: dataCollector,
	}
	analyzer := &analysis.Analyzer{
		Name:     "levels_finalizer",
		Doc:      "report levels check",
		Run:      inst.run,
		Requires: []*analysis.Analyzer{dataCollector.Analyzer}, // ensures that collector will be run first
	}
	inst.Analyzer = analyzer

	return inst
}

func (f *Finalizer) run(pass *analysis.Pass) (interface{}, error) {

	err := processor.ProcessImports(f.commonData)
	if err != nil {
		pass.Reportf(1, err.Error())
	}
	return nil, nil
}
