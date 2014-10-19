package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateXplusYeqZ creates the constraint X + Y = Z
func CreateXplusYeqZ(x, y, z core.VarId) *interval.XplusYeqZ {
	return interval.CreateXplusYeqZ(x, y, z)
}

// CreateXplusYeqC creates the constraint X + Y = C
func CreateXplusYeqC(x, y core.VarId, c int) *interval.C1XplusC2YeqC3 {
	return interval.CreateXplusYeqC(x, y, c)
}

// CreateC1XplusC2YeqC3 creates the constraint C1*X + C2*Y = C3
func CreateC1XplusC2YeqC3(c1 int, x core.VarId, c2 int, y core.VarId,
	c3 int) *interval.C1XplusC2YeqC3 {
	return interval.CreateC1XplusC2YeqC3(c1, x, c2, y, c3)
}

// CreateC1XplusC2YeqC3ZBounds creates the constraint C1*X + C2*Y = C3*Z
// providing bounds consistency.
func CreateC1XplusC2YeqC3ZBounds(c1 int, x core.VarId, c2 int,
	y core.VarId, c3 int, z core.VarId) *interval.C1XplusC2YeqC3ZBounds {
	return interval.CreateC1XplusC2YeqC3ZBounds(c1, x, c2, y, c3, z)
}
