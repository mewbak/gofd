package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/indexical"
)

// b1 + b2 + b3 + b4 + b5 + b6 = 1..#bs

func CreateOr(store *core.Store, constraints []*ReifiedConstraint) core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateOr-propagator")
	}
	varidBools := make([]core.VarId, len(constraints))
	for i, c := range constraints {
		varidBools[i] = c.GetBool()
	}
	resultVar := core.CreateIntVarIvFromTo("OrResult", store, 1, len(constraints))
	return indexical.CreateSum(store, resultVar, varidBools)
}
