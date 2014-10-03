package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateXplusCneqY creates the constraint X + C != Y
func CreateXplusCneqY(x core.VarId, c int, y core.VarId) core.Constraint {
	return interval.CreateXplusCneqY(x, c, y)
}
