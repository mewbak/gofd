package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateSum creates the constraint X1 + X2 + ... + Xn = Z
func CreateSum(store *core.Store,
	resultVar core.VarId, intVars []core.VarId) core.Constraint {
	return interval.CreateSum(store, resultVar, intVars)
}

// CreateSumBounds creates the constraint X1 + X2 + ... + Xn = Z
// providing bounds consistency.
func CreateSumBounds(store *core.Store,
	resultVar core.VarId, intVars []core.VarId) core.Constraint {
	return interval.CreateSumBounds(store, resultVar, intVars)
}

// CreateWeightedSum creates the constraint C1*X1 + C2*X2 + ... + Cn*Xn = Z
func CreateWeightedSum(store *core.Store, resultVar core.VarId, cs []int,
	intVars ...core.VarId) core.Constraint {
	return interval.CreateWeightedSum(store, resultVar, cs, intVars...)
}

// CreateWeightedSumBounds creates the constraint C1*X1 + C2*X2 + ... + Cn*Xn = Z
// providing bounds consistency.
func CreateWeightedSumBounds(store *core.Store, resultVar core.VarId, cs []int,
	intVars ...core.VarId) core.Constraint {
	return interval.CreateWeightedSumBounds(store, resultVar, cs, intVars...)
}
