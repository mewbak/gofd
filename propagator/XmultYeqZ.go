package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateXmultYeqZ creates the constraint X * Y = Z
func CreateXmultYeqZ(x core.VarId, y core.VarId, z core.VarId) *interval.XmultYeqZ {
	return interval.CreateXmultYeqZ(x, y, z)
}

// CreateXmultCeqY creates the constraint X * C = Y
func CreateXmultCeqY(x core.VarId, c int, y core.VarId) *interval.XmultCeqY {
	return interval.CreateXmultCeqY(x, c, y)
}

// CreateC1XmultC2YeqC3ZBounds creates the constraint C1*X * C2*Y = C3*Z
func CreateC1XmultC2YeqC3ZBounds(c1 int, x core.VarId, c2 int,
	y core.VarId, c3 int, z core.VarId) *interval.C1XmultC2YeqC3ZBounds {
	return interval.CreateC1XmultC2YeqC3ZBounds(c1, x, c2, y, c3, z)
}
