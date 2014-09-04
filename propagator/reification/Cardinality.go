package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/indexical"
)

func CreateCardinality(store *core.Store,
	resultVar core.VarId, intVars []core.VarId) core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateCardinality-propagator")
	}
	return indexical.CreateSum(store, resultVar, intVars)
}
