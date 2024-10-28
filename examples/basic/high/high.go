package high

import (
	"fmt"

	"github.com/gennadyterekhov/import-layers-go/examples/basic/low" // want `wrong level. cannot import "github.com/gennadyterekhov/import-layers-go/examples/basic/low" from "github.com/gennadyterekhov/import-layers-go/examples/basic/high"`
)

func Fn() {
	fmt.Println("uplevel")
	low.Fn() // not ok
}
