package main

import (
	"github.com/gennadyterekhov/levelslib/examples/basic/downlevel"
	"github.com/gennadyterekhov/levelslib/examples/basic/uplevel"
)

const LevelslibConfig = "config.json"

// if "mode" : "numbers"
const LevelslibLevel = 1

// if mode : "names"
// const LevelslibLevel = "infra"

func main() {

	uplevel.Fn()   // ok
	downlevel.Fn() // not ok, not compiled
}
