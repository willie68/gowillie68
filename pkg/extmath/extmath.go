package extmath

import (
	"golang.org/x/exp/constraints"
)

// Abs returns the absolute value of x.
func Abs[E constraints.Signed](x E) E {
	if x < 0 {
		return -x
	}
	return x
}
