package explicit

import (
	"bitbucket.org/gofd/gofd/core"
)

// CreateSumBounds creates a sum-constraint with bound consistency.
// Example: X + Y + Q + R = SUM
func CreateSumBounds(store *core.Store, sum core.VarId,
	intVars []core.VarId) []core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateSumBounds")
	}
	len_vars := len(intVars)
	if len_vars == 0 { // empty sum is zero
		return []core.Constraint{CreateXeqC(sum, 0)}
	} else if len_vars == 1 { // just equals
		return []core.Constraint{CreateXeqYBounds(sum, intVars[0])}
	}
	// otherwise build successive sums of two variables
	// e.g. X + Y + Z = SUM iff X + Y = H1, H1 + Z = H2, H2 = SUM
	prop_list := make([]core.Constraint, len_vars)
	H := intVars[0]
	for i, X := range intVars[1:] {
		NewH := core.CreateAuxIntVarExFromTo(store,
			store.GetDomain(H).GetMin()+store.GetDomain(X).GetMin(),
			store.GetDomain(H).GetMax()+store.GetDomain(X).GetMax())
		prop_list[i] = CreateXplusYeqZBounds(H, X, NewH)
		H = NewH
	}
	prop_list[len_vars-1] = CreateXeqYBounds(H, sum)
	return prop_list
}
