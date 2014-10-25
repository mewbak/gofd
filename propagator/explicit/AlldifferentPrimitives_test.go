package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func alldifferentPrimitives_test(t *testing.T, xinit []int, yinit []int,
	zinit []int, qinit []int, expx []int, expy []int,
	expz []int, expq []int, expready bool) {
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	Z := core.CreateIntVarExValues("Z", store, zinit)
	Q := core.CreateIntVarExValues("Q", store, qinit)
	store.AddPropagators(CreateAlldifferentPrimitives(X, Y, Z, Q))
	ready := store.IsConsistent()
	ready_test(t, "Alldifferent2", ready, expready)
	if expready {
		domainEquals_test(t, "Alldifferent2", X, expx)
		domainEquals_test(t, "Alldifferent2", Y, expy)
		domainEquals_test(t, "Alldifferent2", Z, expz)
		domainEquals_test(t, "Alldifferent2", Q, expq)
	}
}

func Test_AlldifferentPrimitivesa(t *testing.T) {
	setup()
	defer teardown()
	log("AlldifferentPrimitivesa: X:0, Y:0..1, Z:1..2, Q:2..3")
	alldifferentPrimitives_test(t, []int{0}, []int{0, 1}, []int{1, 2}, []int{2, 3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_AlldifferentPrimitivesb(t *testing.T) {
	setup()
	defer teardown()
	log("AlldifferentPrimitivesb: X:0..1, Y:1, Z:2..3, Q:3")
	alldifferentPrimitives_test(t, []int{0, 1}, []int{1}, []int{2, 3}, []int{3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_AlldifferentPrimitivesc(t *testing.T) {
	setup()
	defer teardown()
	log("AlldifferentPrimitivesc: X:0, Y:1, Z:2, Q:3")
	alldifferentPrimitives_test(t, []int{0}, []int{1}, []int{2}, []int{3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_AlldifferentPrimitivesd(t *testing.T) {
	setup()
	defer teardown()
	log("AlldifferentPrimitivesd: X:0, Y:0, Z:0, Q:0")
	alldifferentPrimitives_test(t, []int{0}, []int{0}, []int{0}, []int{0},
		[]int{}, []int{}, []int{}, []int{}, false)
}
