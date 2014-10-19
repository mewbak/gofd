package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func alldiff_test(t *testing.T, xinit []int, yinit []int,
	zinit []int, qinit []int, expx []int, expy []int,
	expz []int, expq []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	Y := core.CreateIntVarIvValues("Y", store, yinit)
	Z := core.CreateIntVarIvValues("Z", store, zinit)
	Q := core.CreateIntVarIvValues("Q", store, qinit)
	store.AddPropagators(CreateAlldifferent(X, Y, Z, Q))
	ready := store.IsConsistent()
	ready_test(t, "Alldifferent_IC", ready, expready)
	if expready {
		expX := core.CreateIvDomainFromIntArr(expx)
		expY := core.CreateIvDomainFromIntArr(expy)
		expZ := core.CreateIvDomainFromIntArr(expz)
		expQ := core.CreateIvDomainFromIntArr(expq)

		DomainEquals_test(t, "Alldifferent_IC", X, expX)
		DomainEquals_test(t, "Alldifferent_IC", Y, expY)
		DomainEquals_test(t, "Alldifferent_IC", Z, expZ)
		DomainEquals_test(t, "Alldifferent_IC", Q, expQ)
	}
}

func Test_AllDiff_ICa(t *testing.T) {
	setup()
	defer teardown()
	log("Alldiff_ICa: X:0, Y:0..1, Z:1..2, Q:2..3")
	alldiff_test(t, []int{0}, []int{0, 1}, []int{1, 2}, []int{2, 3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_AllDiff_ICb(t *testing.T) {
	setup()
	defer teardown()
	log("Alldiff_ICb: X:0..1, Y:1, Z:2..3, Q:3")
	alldiff_test(t, []int{0, 1}, []int{1}, []int{2, 3}, []int{3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_AllDiff_ICc(t *testing.T) {
	setup()
	defer teardown()
	log("Alldiff_ICc: X:0, Y:1, Z:2, Q:3")
	alldiff_test(t, []int{0}, []int{1}, []int{2}, []int{3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_AllDiff_ICd(t *testing.T) {
	setup()
	defer teardown()
	log("Alldiff_ICd: X:0, Y:0, Z:0, Q:0")
	alldiff_test(t, []int{0}, []int{0}, []int{0}, []int{0},
		[]int{}, []int{}, []int{}, []int{}, false)
}
