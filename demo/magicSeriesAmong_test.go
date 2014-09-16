package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"fmt"
	"testing"
)

// magic series problem modelled with among
// sequence of variables X0,...,Xn
// variable X0 declares how often 0 can occur in the sequence of assignments
// variable X1 declares how often 1 can occur etc.

// testMagicSeries creates a magic series problem with n variables
func testMagicSeries(t *testing.T, n int, expectedResult bool) {
	log(fmt.Sprintf("magic series with    Among: n = %d", n))
	// define variables X0,...,Xn
	variables := make([]core.VarId, n+1)
	for i := 0; i < len(variables); i++ {
		variables[i] = core.CreateAuxIntVarExFromTo(store, 0, n)
	}
	// define constraints
	// each value j can occur Xj times
	for i := 0; i < len(variables); i++ {
		store.AddPropagator(explicit.CreateAmong(variables,
			[]int{i}, variables[i]))
	}
	query := labeling.CreateSearchOneQueryVariableSelect(variables)
	// labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	labeling.Labeling(store, query,
		labeling.SmallestDomainFirst, labeling.InDomainMin)
	ready := store.IsConsistent()
	log(fmt.Sprintf("ready: %6v,  search nodes=%4d",
		ready, query.GetSearchStatistics().GetNodes()))
	ready_test(t, "Magic Series", ready, expectedResult)
}

func Test_magicSeries6(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeries(t, 6, true)
}

func Test_magicSeries7(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeries(t, 7, true)
}

func Test_magicSeries8(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeries(t, 8, true)
}

func Test_magicSeries9(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeries(t, 9, true)
}

func Test_magicSeries10(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeries(t, 10, true)
}
