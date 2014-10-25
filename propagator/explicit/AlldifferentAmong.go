package explicit

import (
	"bitbucket.org/gofd/gofd/core"
)

// implementation of Alldifferent Constraint with Among
// signature: Alldiff({X1,...,Xi})
// during the assignment of Xi no value can occur twice
// CreateAlldifferentAmong creates an alldifferent constraint modelled
// with Among.
func CreateAlldifferentAmong(xi []core.VarId,
	store *core.Store) []core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateAlldifferentAmong_propagator")
	}
	// create the union of all domains of Xi
	domain := createUnionFromDomains(xi, store)
	// make an Among constraint for every value in the union of
	// the domains of Xi and set D(N) to {0,1}
	amongPropagators := make([]core.Constraint, len(domain))
	for i := 0; i < len(domain); i++ {
		amongPropagators[i] = CreateAmong(xi, []int{domain[i]},
			core.CreateAuxIntVarExFromTo(store, 0, 1))
	}
	return amongPropagators
}
