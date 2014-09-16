package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateAmong creates an Among constraint forcing that n variables of
// the xi variables take values from the ones in k.
func CreateAmong(xi []core.VarId, k []int, n core.VarId) *interval.Among {
	return interval.CreateAmong(xi, k, n)
}
