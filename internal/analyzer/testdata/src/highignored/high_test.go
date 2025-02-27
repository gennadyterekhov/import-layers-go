package highignored

import (
	"testing"

	"low" // ok because ignored
)

func FnTest() {
	low.Fn()
}

func TestEq(t *testing.T) {

}
