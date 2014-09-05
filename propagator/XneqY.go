package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateXneqY creates the constraint X != Y
func CreateXneqY(x core.VarId, y core.VarId) *interval.XneqY {
	return interval.CreateXneqY(x, y)
}

// CreateXneqC creates the constraint X != C
func CreateXneqC(x core.VarId, c int) *interval.XneqC {
	return interval.CreateXneqC(x, c)
}
