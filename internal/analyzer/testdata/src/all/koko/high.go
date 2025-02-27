package high

import (
	"low" // want `cannot import package from lower layer`
)

func Fn() {
	low.Fn()
}
