package highmid

import (
	"low" // want `cannot import package from lower layer`

	"mid" // want `cannot import package from lower layer`
)

func Fn() {
	low.Fn()
	mid.Mid()
}
