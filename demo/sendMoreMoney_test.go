package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"testing"
)

func Test_sendMoreMoneyPlain(t *testing.T) {
	setup()
	defer teardown()
	log("sendMoreMoneyPlain      : E,N,D,O,R,Y: 0..9; S,M: 1..9")
	doSMM(t, false)
}

func Test_sendMoreMoneyIndexicals(t *testing.T) {
	setup()
	defer teardown()
	log("sendMoreMoneyIndexicals : E,N,D,O,R,Y: 0..9; S,M: 1..9")
	doSMM(t, true)
}

func doSMM(t *testing.T, withIndexical bool) {
	varIds := ConstrainSendMoreMoney(store, withIndexical)
	labeling.SetAllvars(varIds)
	query := labeling.CreateSearchOneQuery()
	labeling.Labeling(store, query, labeling.InDomainMin, labeling.VarSelect)
	domainCheck(t, "sendMoreMoney", query, varIds...)
	searchStat(query.GetSearchStatistics())
}

func unpack(varIda []core.VarId, varIds ...*core.VarId) {
	for i, varId := range varIda {
		*varIds[i] = varId
	}
}

func domainCheck(t *testing.T,
	test_id string, query *labeling.SearchOneQuery,
	vars ...core.VarId) {
	resultSet := query.GetResultSet()
	if len(resultSet) != 1 {
		t.Errorf("%s nosols = %d, want %d",
			test_id, len(resultSet), 1)
	}
	var S, E, N, D, M, O, R, Y core.VarId
	unpack(vars, &S, &E, &N, &D, &M, &O, &R, &Y)
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
