// Package data  used only in tests
package data

import (
	"sync"
)

type CommonData struct {
	pakToLev map[string]int
	deptree  map[string][]string
	lock     sync.Mutex
}

func New() *CommonData {
	return &CommonData{
		pakToLev: make(map[string]int, 0),
		deptree:  make(map[string][]string, 0),
	}
}

func (cd *CommonData) AddImport(pkgPath, imported string) {
	cd.lock.Lock()
	cd.deptree[pkgPath] = append(cd.deptree[pkgPath], imported)
	cd.lock.Unlock()
}

func (cd *CommonData) AddPackage(pkgPath string, level int) {
	cd.lock.Lock()
	cd.pakToLev[pkgPath] = level
	cd.lock.Unlock()
}

func (cd *CommonData) PakToLev() map[string]int {
	return cd.pakToLev
}

func (cd *CommonData) Deptree() map[string][]string {
	return cd.deptree
}
