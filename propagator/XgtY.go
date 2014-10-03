package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateXgtYplusC creates the constraint X > Y+C
func CreateXgtYplusC(x core.VarId, y core.VarId, c int) *interval.XgtYplusC {
	return interval.CreateXgtYplusC(x, y, c)
}

// CreateXgteqY creates the constraint X >= Y
func CreateXgteqY(x core.VarId, y core.VarId) *interval.XgtYplusC {
	return interval.CreateXgteqY(x, y)
}

// CreateXgtY creates the constraint X > Y
func CreateXgtY(x core.VarId, y core.VarId) *interval.XgtY {
	return interval.CreateXgtY(x, y)
}

// CreateXgtC creates the constraint X > C
func CreateXgtC(x core.VarId, c int) *interval.XgtC {
	return interval.CreateXgtC(x, c)
}
