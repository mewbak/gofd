package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"testing"
)

// X=5<=>B0, Y=10<=>B1, Cardinality([B0,B1],d)
// Same as simplereification test

func Test_Cardinality1(t *testing.T) {
	setup()
	defer teardown()
	//core.GetLogger().SetLoggingLevel(core.LOG_DEBUG)
	log("Cardinality1: X:1..10, Y:1..10, X=5<=>B0, Y=10<=>B1, B0+B1=2")

	doCardinalityTest(t, 1, 10, 5, 1, 10, 10, [][]int{{5, 5}}, [][]int{{10, 10}}, 2, 2, true)
}

func Test_Cardinality2(t *testing.T) {
	setup()
	defer teardown()
	log("Cardinality2: X:1..10, Y:1..10, X=5<=>B0, Y=10<=>B1, B0+B1=0")

	doCardinalityTest(t, 1, 10, 5, 1, 10, 10, [][]int{{1, 4}, {6, 10}}, [][]int{{1, 9}}, 0, 0, true)
}

//delayed
func Test_Cardinality3(t *testing.T) {
	setup()
	defer teardown()
	log("Cardinality3: X:1..10, Y:1..10, X=5<=>B0, Y=10<=>B1, B0+B1=1")

	doCardinalityTest(t, 1, 10, 5, 1, 10, 10, [][]int{{1, 10}}, [][]int{{1, 10}}, 1, 1, true)
}

func Test_Cardinality4(t *testing.T) {
	setup()
	defer teardown()
	log("Cardinality4: X:1..10, Y:1..10, X=5<=>B0, Y=10<=>B1, B0+B1=0")

	doCardinalityTest(t, 1, 10, 5, 1, 10, 10, [][]int{{1, 10}}, [][]int{{1, 10}}, 0, 2, true)
}

func doCardinalityTest(t *testing.T,
	fromX int, toX int, eqX int,
	fromY int, toY int, eqY int,
	expX [][]int,
	expY [][]int,
	trueConstraintsMin int, trueConstraintsMax int,
	expResult bool) {

	var X, Y, B0, B1, B core.VarId

	core.CreateIntVarsIvFromTo([]*core.VarId{&X},
		[]string{"X"}, store, fromY, toY)
	core.CreateIntVarsIvFromTo([]*core.VarId{&X},
		[]string{"Y"}, store, fromY, toY)

	core.CreateIntVarsIvFromTo([]*core.VarId{&B0, &B1},
		[]string{"B0", "B1"}, store, 0, 1)
	core.CreateIntVarsIvFromTo([]*core.VarId{&B},
		[]string{"B"}, store, trueConstraintsMin, trueConstraintsMax)

	xeqc := indexical.CreateXeqC(X, eqX)
	yeqc := indexical.CreateXeqC(Y, eqY)

	ricXeqC := CreateReifiedConstraint(xeqc, B0)
	ricYeqC := CreateReifiedConstraint(yeqc, B1)

	store.AddPropagator(ricXeqC)
	store.AddPropagator(ricYeqC)
	store.AddPropagator(CreateCardinality(store, B, []core.VarId{B0, B1}))

	ready := store.IsConsistent()
	ready_test(t, "Cardinality", ready, expResult)

	if ready == expResult {
		expXDom := core.CreateIvDomainFromTos(expX)
		expYDom := core.CreateIvDomainFromTos(expY)
		DomainEquals_test(t, "Cardinality", X, expXDom)
		DomainEquals_test(t, "Cardinality", Y, expYDom)
	}

	propStat()
}
