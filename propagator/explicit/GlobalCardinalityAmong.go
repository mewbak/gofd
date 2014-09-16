package explicit

import (
	"bitbucket.org/gofd/gofd/core"
)

// implementation of Global Cardinality constraint with Among
// signature: GCC({X1,...,Xi}, k, counters)
// every value k[j] appears exactly counters[j] times while assigning
// values to Xi

// CreateGCCAmong creates a Global Cardinality constraint modelled with Among
func CreateGCCAmong(xi []core.VarId, k []int, counters []int,
	store *core.Store) []core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateGCC_propagator")
	}

	// the amoung of counters and the amount of values in k have to be identical
	amountCounters := len(counters)
	if len(k) != amountCounters {
		panic("The length of counters is not equal to the amount " +
			"of elements in the domain of k!!")
		return nil
	}
	// for every value k[i] an Among constraint is created which limits
	// the number of occurrences of k[i] in the assignment to
	// exactly counters[i] times
	amongPropagators := make([]core.Constraint, amountCounters)
	for i := 0; i < len(k); i++ {
		amongPropagators[i] = CreateAmong(xi, []int{k[i]},
			core.CreateAuxIntVarExFromTo(store, counters[i], counters[i]))
	}
	return amongPropagators
}
