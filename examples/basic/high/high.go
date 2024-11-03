package high

import (
	"fmt"

	"github.com/gennadyterekhov/import-layers-go/examples/basic/low" // want `cannot import package from lower layer`
)

func Fn() {
	fmt.Println("high")
	low.Fn() // not ok
}
