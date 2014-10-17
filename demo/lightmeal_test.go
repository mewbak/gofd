package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"testing"
)

// Einführung in die Constraint-Programmierung; Hofstedt, Wolf
// A meal consists of an appetizer, a main dish and a desert. Each
// dish has a weight, the higher the heavier the dish (e.g. kcal).
// As appetizers there is pasta with 4 and radishes with 1. As main
// there is pork with 11 and beef with 7 as desert there is fruit
// with 2 and icecream with 6. A light meal is a meal with a
// cumulated weight of 10.
func Test_lightmeal(t *testing.T) {
	setup()
	defer teardown()
	log("lightmeal")
	appetizer := core.CreateIntVarIvValues("appetizer", store, []int{1, 4})
	main := core.CreateIntVarIvValues("main", store, []int{11, 7})
	desert := core.CreateIntVarIvValues("desert", store, []int{2, 6})
	sum := core.CreateIntVarIvValues("sum", store, []int{10})
	store.AddPropagators(propagator.CreateSum(store, sum,
		[]core.VarId{appetizer, main, desert}))
	ready := store.IsConsistent()
	ready_test(t, "lightmeal", ready, true)
	domainEquals_test(t, "lightmeal", appetizer, []int{1})
	domainEquals_test(t, "lightmeal", main, []int{7})
	domainEquals_test(t, "lightmeal", desert, []int{2})
	propStat()
}
