// Package data  used only in tests
package data

import (
	"sync"
)

type CommonData struct {
	packageToLayer map[string]int
	deptree        map[string][]string
	lock           sync.Mutex
}

func New() *CommonData {
	return &CommonData{
		packageToLayer: make(map[string]int, 0),
		deptree:        make(map[string][]string, 0),
	}
}

func (cd *CommonData) AddImport(pkgPath, imported string) {
	cd.lock.Lock()
	cd.deptree[pkgPath] = append(cd.deptree[pkgPath], imported)
	cd.lock.Unlock()
}

func (cd *CommonData) AddPackage(pkgPath string, layer int) {
	cd.lock.Lock()
	cd.packageToLayer[pkgPath] = layer
	cd.lock.Unlock()
}

func (cd *CommonData) PackageToLayer() map[string]int {
	return cd.packageToLayer
}

func (cd *CommonData) Deptree() map[string][]string {
	return cd.deptree
}
