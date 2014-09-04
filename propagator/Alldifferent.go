package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateAlldifferent creates a constraint for pairwise difference,
// e.g. X≠Y, X≠Z, Y≠Z for X,Y,Z. It behaves like quadratically many
// "not equal" constraints.
func CreateAlldifferent(vars ...core.VarId) core.Constraint {
	return interval.CreateAlldifferent(vars...)
}

// CreateAlldifferent_Offset creates a constraint similar to Alldifferent,
// but allows to specify an offset for each variable. Thus,
// e.g. X+dX≠Y+dY, X+dX≠Z+dZ, Y+dY≠Z+dZ must hold for three variables
// {X, Y, Z} and offsets {dX, dY, dZ}.
func CreateAlldifferent_Offset(vars []core.VarId, offsets []int) core.Constraint {
	return interval.CreateAlldifferent_Offset(vars, offsets)
}
