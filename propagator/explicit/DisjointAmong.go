package explicit

import (
	"bitbucket.org/gofd/gofd/core"
)

//implementation of the Disjoint constraint with Among
//signature: Disjoint({X1,...,Xi}, {Y1,...,Yj})
//no value that appears in Xi can appear in Yj, but Xi and Yj are multisets

//CreateDisjointAmong creates a Disjoint constraint modelled with Among
func CreateDisjointAmong(xi []core.VarId, yj []core.VarId, store *core.Store) []core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateDisjoint_propagator")
	}

	//make unions of the domains so the unions can be used as K
	xiUnion := createUnionFromDomains(xi, store)
	yjUnion := createUnionFromDomains(yj, store)

	return []core.Constraint{CreateAmong(xi, yjUnion, core.CreateAuxIntVarExFromTo(store, 0, 0)),
		CreateAmong(yj, xiUnion, core.CreateAuxIntVarExFromTo(store, 0, 0))}
}
