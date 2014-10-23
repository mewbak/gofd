package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"bitbucket.org/gofd/gofd/propagator/indexical"
)

//     S E N D         9 5 6 7
// +   M O R E     +   1 0 8 5
// -----------     -----------
// = M O N E Y     = 1 0 6 5 2

// creates the SEND+MORE=MONEY problem and returns the
// created variables [S,E,N,D,M,O,R,Y]
func ConstrainSendMoreMoney(store *core.Store,
	withIndexical bool) []core.VarId {
	S := core.CreateIntVarFromTo("S", store, 1, 9)
	E := core.CreateIntVarFromTo("E", store, 0, 9)
	N := core.CreateIntVarFromTo("N", store, 0, 9)
	D := core.CreateIntVarFromTo("D", store, 0, 9)
	M := core.CreateIntVarFromTo("M", store, 1, 9)
	O := core.CreateIntVarFromTo("O", store, 0, 9)
	R := core.CreateIntVarFromTo("R", store, 0, 9)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	var alldiff_prop core.Constraint
	if !withIndexical {
		alldiff_prop = propagator.CreateAlldifferent(S, E, N, D, M, O, R, Y)
	} else {
		alldiff_prop = indexical.CreateAlldifferent(S, E, N, D, M, O, R, Y)
	}
	store.AddPropagators(alldiff_prop)
	max := 1000*store.GetDomain(S).GetMax() + 100*store.GetDomain(E).GetMax() +
		10*store.GetDomain(N).GetMax() + 1*store.GetDomain(D).GetMax() +
		1000*store.GetDomain(M).GetMax() + 100*store.GetDomain(O).GetMax() +
		10*store.GetDomain(R).GetMax() + 1*store.GetDomain(E).GetMax()
	sum := core.CreateIntVarFromTo("sum", store, 0, max)
	weights := []int{1000, 100, 10, 1, 1000, 100, 10, 1}
	vars := []core.VarId{S, E, N, D, M, O, R, E}
	var sendmore_props core.Constraint
	if !withIndexical {
		sendmore_props = propagator.CreateWeightedSumBounds(store,
			sum, weights, vars...)
	} else {
		sendmore_props = indexical.CreateWeightedSumBounds(store,
			sum, weights, vars...)
	}
	store.AddPropagators(sendmore_props)
	weights = []int{10000, 1000, 100, 10, 1}
	vars = []core.VarId{M, O, N, E, Y}
	var money_props core.Constraint
	if !withIndexical {
		money_props = propagator.CreateWeightedSumBounds(store,
			sum, weights, vars...)
	} else {
		money_props = indexical.CreateWeightedSumBounds(store,
			sum, weights, vars...)
	}
	store.AddPropagators(money_props)
	return []core.VarId{S, E, N, D, M, O, R, Y}
}
