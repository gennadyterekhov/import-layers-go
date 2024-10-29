package high

import (
	"low" // want `wrong layer. cannot import "low" from "high"`
)

func Fn() {
	low.Fn()
}
