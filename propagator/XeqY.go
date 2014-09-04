package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateXeqC creates the constraint X = C
func CreateXeqC(x core.VarId, c int) *interval.XeqC {
	return interval.CreateXeqC(x, c)
}

// CreateXplusCeqY creates the constraint X+C = Y
func CreateXplusCeqY(x core.VarId, c int, y core.VarId) *interval.XplusCeqY {
	return interval.CreateXplusCeqY(x, c, y)
}

// CreateC1XeqC2YBounds creates the constraint C1*X = C2*Y
// providing bounds consistency.
func CreateC1XeqC2YBounds(c1 int, x core.VarId,
	c2 int, y core.VarId) *interval.C1XeqC2YBounds {
	return interval.CreateC1XeqC2YBounds(c1, x, c2, y)
}
