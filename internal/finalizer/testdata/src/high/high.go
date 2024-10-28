package high

import (
	"low" // want `wrong level. cannot import "low" from "high"`
)

func Fn() {
	low.Fn()
}
