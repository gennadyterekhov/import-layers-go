package high

import (
	"fmt"

	"github.com/gennadyterekhov/import-layers-go/examples/basic/low" // want `wrong layer. cannot import "github.com/gennadyterekhov/import-layers-go/examples/basic/low" from "github.com/gennadyterekhov/import-layers-go/examples/basic/high"`
)

func Fn() {
	fmt.Println("high")
	low.Fn() // not ok
}
