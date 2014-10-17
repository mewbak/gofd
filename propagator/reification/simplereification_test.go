package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"testing"
)

// X=5<=>B0, Y=10<=>B1, B0+B1=d

func Test_simplereificationExample1(t *testing.T) {
	setup()
	defer teardown()
	//core.GetLogger().SetLoggingLevel(core.LOG_DEBUG)
	log("simplereification1: X:1..10, Y:1..10, X=5<=>B0, Y=10<=>B1, B0+B1=2")

	DoSimpleReification(t, 1, 10, 5, 1, 10, 10, [][]int{{5, 5}}, [][]int{{10, 10}}, 2, 2, true)
}

func Test_simplereificationExample2(t *testing.T) {
	setup()
	defer teardown()
	log("simplereification2: X:1..10, Y:1..10, X=5<=>B0, Y=10<=>B1, B0+B1=0")

	DoSimpleReification(t, 1, 10, 5, 1, 10, 10, [][]int{{1, 4}, {6, 10}}, [][]int{{1, 9}}, 0, 0, true)
}

//delayed
func Test_simplereificationExample3(t *testing.T) {
	setup()
	defer teardown()
	log("simplereification3: X:1..10, Y:1..10, X=5<=>B0, Y=10<=>B1, B0+B1=1")

	DoSimpleReification(t, 1, 10, 5, 1, 10, 10, [][]int{{1, 10}}, [][]int{{1, 10}}, 1, 1, true)
}

func DoSimpleReification(t *testing.T,
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
	store.AddPropagator(indexical.CreateSum(store, B, []core.VarId{B0, B1}))

	ready := store.IsConsistent()
	ready_test(t, "simplereification", ready, expResult)

	if ready == expResult {
		expXDom := core.CreateIvDomainFromTos(expX)
		expYDom := core.CreateIvDomainFromTos(expY)
		DomainEquals_test(t, "simplereification", X, expXDom)
		DomainEquals_test(t, "simplereification", Y, expYDom)
	}

}
