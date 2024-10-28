package main

import (
	"github.com/gennadyterekhov/import-layers-go/examples/basic/high"
	"github.com/gennadyterekhov/import-layers-go/examples/basic/low"
)

func main() {
	high.Fn() // ok
	low.Fn()  // ok
}
