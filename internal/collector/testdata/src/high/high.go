package high

import (
	"low" // want `wrong level. cannot import low from high`
)

const LevelslibLevel = 1

func Fn() {
	low.Fn()
}
