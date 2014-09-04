package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"bitbucket.org/gofd/gofd/propagator/reification"
	"fmt"
	"testing"
)

// magic series problem modelled without among
// sequence of variables X0,...,Xn
// variable X0 declares how often 0 can occur in the sequence of variables
// variable X1 declares how often 1 can occur etc.

// testMagicSeriesWithoutAmong creates a magic series problem with n variables
func testMagicSeriesWithoutAmong(t *testing.T, n int, expectedResult bool) {
	log(fmt.Sprintf("magic series without Among: n = %d", n))
	// define variables X0,...,Xn
	variables := make([]core.VarId, n+1)
	for i := 0; i < len(variables); i++ {
		variables[i] = core.CreateAuxIntVarIvFromTo(store, 0, n)
	}
	// define constraints
	// each value j can occur Xj times
	for i := 0; i < len(variables); i++ {
		// array for reified constraints
		variables := make([]core.VarId, len(variables))
		for j := 0; j < len(variables); j++ {
			// store in variables[j] whether Xj (variables[j]) takes
			// the value i or not
			variables[j] = core.CreateAuxIntVarIvFromTo(store, 0, 1)
			xeqc := indexical.CreateXeqC(variables[j], i)
			reifiedConstraint := reification.CreateReifiedConstraint(xeqc,
				variables[j])
			store.AddPropagator(reifiedConstraint)
		}
		// the amount of variables in X0,...,Xn that have taken the value i
		// must correspond to Xi (variables[i])
		store.AddPropagator(interval.CreateSum(store, variables[i], variables))
	}
	query := labeling.CreateSearchOneQueryVariableSelect(variables)
	labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	ready := store.IsConsistent()
	log(fmt.Sprintf("ready: %6v,  search nodes=%4d",
		ready, query.GetSearchStatistics().GetNodes()))
	ready_test(t, "Magic Series Without Among", ready, expectedResult)
}

func Test_magicSeriesWithoutAmong6(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeriesWithoutAmong(t, 6, true)
}

func Test_magicSeriesWithoutAmong7(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeriesWithoutAmong(t, 7, true)
}

func Test_magicSeriesWithoutAmong8(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeriesWithoutAmong(t, 8, true)
}

func Test_magicSeriesWithoutAmong9(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeriesWithoutAmong(t, 9, true)
}

func Test_magicSeriesWithoutAmong10(t *testing.T) {
	setup()
	defer teardown()
	testMagicSeriesWithoutAmong(t, 10, true)
}
