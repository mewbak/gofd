package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"testing"
)

//     S E N D         9 5 6 7
// +   M O R E     +   1 0 8 5
// -----------     -----------
// = M O N E Y     = 1 0 6 5 2

func Test_smm(t *testing.T) {
	setup()
	defer teardown()
	log("send_more_money                 : E,N,D,O,R,Y: 0..9; S,M: 1..9")
	doSMM(t, false)
}

func Test_smm_indexicals(t *testing.T) {
	setup()
	defer teardown()
	log("send_more_money with indexicals : E,N,D,O,R,Y: 0..9; S,M: 1..9")
	doSMM(t, true)
}

// interval domains and bounds consistency, optional with indexicals
func doSMM(t *testing.T, withIndexical bool) {
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
	labeling.SetAllvars([]core.VarId{S, E, N, D, M, O, R, Y})
	query := labeling.CreateSearchOneQuery()
	labeling.Labeling(store, query, labeling.InDomainMin, labeling.VarSelect)
	domainCheck(t, "send more money", query, S, E, N, D, M, O, R, Y)
	searchStat(query.GetSearchStatistics())
}

func domainCheck(t *testing.T,
	test_id string, query *labeling.SearchOneQuery,
	S, E, N, D, M, O, R, Y core.VarId) {
	resultSet := query.GetResultSet()
	if len(resultSet) != 1 {
		t.Errorf("%s nosols = %d, want %d",
			test_id, len(resultSet), 1)
	}
	valueEqual(t, test_id, S, resultSet[0][S], 9)
	valueEqual(t, test_id, E, resultSet[0][E], 5)
	valueEqual(t, test_id, N, resultSet[0][N], 6)
	valueEqual(t, test_id, D, resultSet[0][D], 7)
	valueEqual(t, test_id, M, resultSet[0][M], 1)
	valueEqual(t, test_id, O, resultSet[0][O], 0)
	valueEqual(t, test_id, R, resultSet[0][R], 8)
	valueEqual(t, test_id, Y, resultSet[0][Y], 2)
}

func valueEqual(t *testing.T, test_id string,
	v_id core.VarId, calcVal int, expVal int) {
	if calcVal != expVal {
		t.Errorf("%s, for var_id %d calculated domain = %d, want %d",
			test_id, v_id, calcVal, expVal)
	}
}
