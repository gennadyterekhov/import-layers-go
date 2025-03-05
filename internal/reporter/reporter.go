package reporter

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Report struct {
	pos token.Pos
	str string
}

type Reporter struct {
	Report     func(pos token.Pos, format string, args ...interface{})
	reports    []Report
	reportsMap map[string][]string
}

func New(report func(pos token.Pos, format string, args ...interface{})) *Reporter {
	return &Reporter{
		Report:     report,
		reports:    make([]Report, 0),
		reportsMap: make(map[string][]string),
	}
}

func NewMock() *Reporter {
	inst := &Reporter{
		reports:    make([]Report, 0),
		reportsMap: make(map[string][]string),
	}
	inst.Report = inst.mockReport

	return inst
}

func (r *Reporter) GetReports() []Report {
	return r.reports
}

func (r *Reporter) AddReport(faultyNode *ast.File, importNode *ast.ImportSpec) {
	//r.beforeReport()
	//importedPkgPath := strings.Trim(importNode.Path.Value, "\"")
	//nodeName := faultyNode.Name.String()
	//
	//v, ok := r.reportsMap[nodeName]
	//if ok && slices.Contains(v, importedPkgPath) {
	//	//return
	//}
	//
	//if r.reportsMap[nodeName] == nil {
	//	r.reportsMap[nodeName] = make([]string, 0)
	//}
	//r.reportsMap[nodeName] = append(r.reportsMap[nodeName], importedPkgPath)

	r.Report(importNode.Pos(), "cannot import package from lower layer")
	//r.afterReport()
}

func (r *Reporter) beforeReport() {

}

func (r *Reporter) afterReport() {

}

func (r *Reporter) mockReport(pos token.Pos, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	r.reports = append(r.reports, Report{pos: pos, str: msg})
}
