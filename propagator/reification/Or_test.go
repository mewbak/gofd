package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"testing"
)

func or_test(t *testing.T, constraints []*ReifiedConstraint, expready bool) {
	or := CreateOr(store, constraints)
	for _, c := range constraints {
		store.AddPropagator(c)
	}
	store.AddPropagator(or)
	ready := store.IsConsistent()
	ready_test(t, "Or", ready, expready)
}

func Test_Or_a(t *testing.T) {
	setup()
	defer teardown()
	log("Or_b")
	// X:0...10
	// X=5
	// X!=5
	X := core.CreateIntVarIvFromTo("X", store, 0, 10)
	B1 := core.CreateIntVarIvDomBool("B1", store)
	B2 := core.CreateIntVarIvDomBool("B2", store)
	// X=5
	xeq := indexical.CreateXeqC(X, 5)
	// X!=5
	xneq := indexical.CreateXneqC(X, 5)
	rc1 := CreateReifiedConstraint(xeq, B1)
	rc2 := CreateReifiedConstraint(xneq, B2)
	or_test(t, []*ReifiedConstraint{rc1, rc2}, true)
}

func Test_Or_b(t *testing.T) {
	setup()
	defer teardown()
	log("Or_c")
	// X:0...10
	// X=12
	X := core.CreateIntVarIvFromTo("X", store, 0, 10)
	B1 := core.CreateIntVarIvDomBool("B1", store)
	// X=12
	xeq := indexical.CreateXeqC(X, 12)
	rc1 := CreateReifiedConstraint(xeq, B1)
	or_test(t, []*ReifiedConstraint{rc1}, false)
}
