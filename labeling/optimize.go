package labeling

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"math"
)

// Maximize maximizes the value of a variable valueVarId (the objective
// function) and returns the maximal value as well as one witness of
// that solution (an assignemtn of all the variables in the store)
// The store is unchanged after the operation.
func Maximize(store *core.Store,
	valueVarId core.VarId) (maxValue int, solution map[core.VarId]int) {
	maxValue, solution, _, _ = maximize(store, valueVarId, false)
	return
}

// MaximizeStates maximizes and returns search and store statistics
func MaximizeStats(store *core.Store,
	valueVarId core.VarId) (maxValue int, solution map[core.VarId]int,
	searchStat *SearchStatistics, storeStat *core.StoreStatistics) {
	maxValue, solution, searchStat, storeStat = maximize(store, valueVarId, true)
	return
}

// Maximize maximizes the value of a variable with respect to a store
func maximize(store *core.Store, valueVar core.VarId,
	doStats bool) (int, map[core.VarId]int, // maximal value and assignment
	*SearchStatistics, *core.StoreStatistics) {
	result := true
	curValue := math.MinInt32
	var searchStats *SearchStatistics = nil
	if doStats {
		searchStats = CreateSearchStatistics()
	}
	store.IsConsistent() // should better be
	var lastResult map[core.VarId]int
	for result { // increase bound incrementally until unsatisfiable
		oStore := store.Clone(nil)
		oStore.AddPropagator(propagator.CreateXgtC(valueVar, curValue))
		query := CreateSearchOneQuery()
		result = Labeling(oStore, query)
		if doStats {
			stats := query.GetSearchStatistics()
			searchStats.UpdateSearchStatistics(stats)
		}
		if result {
			lastResult = query.GetResultSet()[0]
			curValue = lastResult[valueVar]
			// fmt.Printf("  solution value: %d\n", curValue)
		}
	}
	var storeStats *core.StoreStatistics = nil
	if doStats {
		storeStats = searchStats.GetStoreStatistics()
	}
	return curValue, lastResult, searchStats, storeStats
}
