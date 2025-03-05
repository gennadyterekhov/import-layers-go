package highlogicandtest

import (
	"low" // want `cannot import package from lower layer`
)

func Fn() {
	low.Fn()
}
