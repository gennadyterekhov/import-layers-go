package hightest

import (
	"testing"

	"low" // want `cannot import package from lower layer`
)

func FnTest() {
	low.Fn()
}

func TestEq(t *testing.T) {

}
