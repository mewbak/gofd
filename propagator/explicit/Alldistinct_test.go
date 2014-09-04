package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func alldistinct_test(t *testing.T, xinit []int, yinit []int,
	zinit []int, qinit []int, expx []int, expy []int,
	expz []int, expq []int, expready bool) []core.VarId {
	X, Y, Z, Q := createVars(xinit, yinit, zinit, qinit)
	varIds := make([]core.VarId, 4)
	varIds[0], varIds[1], varIds[2], varIds[3] = X, Y, Z, Q
	store.AddPropagators(CreateAlldistinct(X, Y, Z, Q))
	ready := store.IsConsistent()
	ready_test(t, "Alldistinct", ready, expready)
	if expready {
		domainEquals_test(t, "Alldistinct", X, expx)
		domainEquals_test(t, "Alldistinct", Y, expy)
		domainEquals_test(t, "Alldistinct", Z, expz)
		domainEquals_test(t, "Alldistinct", Q, expq)
	}
	return varIds
}

func createVars(xinit []int, yinit []int,
	zinit []int, qinit []int) (core.VarId, core.VarId, core.VarId, core.VarId) {
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	Z := core.CreateIntVarExValues("Z", store, zinit)
	Q := core.CreateIntVarExValues("Q", store, qinit)
	return X, Y, Z, Q
}

func Test_Alldistincta(t *testing.T) {
	setup()
	defer teardown()
	log("Alldistincta: X:0, Y:0..1, Z:1..2, Q:2..3")
	alldistinct_test(t, []int{0}, []int{0, 1}, []int{1, 2}, []int{2, 3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_Alldistinctb(t *testing.T) {
	setup()
	defer teardown()
	log("Alldistinctb: X:0..1, Y:1, Z:2..3, Q:3")
	alldistinct_test(t, []int{0, 1}, []int{1}, []int{2, 3}, []int{3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_Alldistinctc(t *testing.T) {
	setup()
	defer teardown()
	log("Alldistinctc: X:0, Y:1, Z:2, Q:3")
	alldistinct_test(t, []int{0}, []int{1}, []int{2}, []int{3},
		[]int{0}, []int{1}, []int{2}, []int{3}, true)
}

func Test_Alldistinctd(t *testing.T) {
	setup()
	defer teardown()
	log("Alldistinctd: X:0, Y:0, Z:0, Q:0")
	alldistinct_test(t, []int{0}, []int{0}, []int{0}, []int{0},
		[]int{}, []int{}, []int{}, []int{}, false)
}

func Test_Alldistincte(t *testing.T) {
	setup()
	defer teardown()
	log("Alldistincte: X:0..2, Y:0..2, Z:0..2, Q:0..2")
	alldistinct_test(t,
		[]int{0, 1, 2}, []int{0, 1, 2}, []int{0, 1, 2}, []int{0, 1, 2},
		[]int{}, []int{}, []int{}, []int{}, false)
}

// Alldistinct with XltC
func Test_Alldistinctf(t *testing.T) {
	setup()
	defer teardown()
	log("Alldistinctf: X:0..2, Y:0..2, Z:0..2, Q:0..3")
	varIds := alldistinct_test(t,
		[]int{0, 1, 2}, []int{0, 1, 2}, []int{0, 1, 2}, []int{0, 1, 2, 3},
		[]int{0, 1, 2}, []int{0, 1, 2}, []int{0, 1, 2}, []int{3}, true)
	store.AddPropagator(CreateXltC(varIds[3], 3)) // remove 3 from Q
	ready := store.IsConsistent()
	ready_test(t, "Alldistinct", ready, false) // shall fail
}

func Test_Alldistinctg(t *testing.T) {
	setup()
	defer teardown()
	log("Alldistinctg: X:0..1, Y:0..1, Z:0..1, Q:1..3")
	alldistinct_test(t,
		[]int{0, 1}, []int{0, 1}, []int{0, 1}, []int{1, 2, 3},
		nil, nil, nil, nil, false)
}

func Test_Alldistincth(t *testing.T) {
	setup()
	defer teardown()
	log("Alldistincth: X:0..4, Y:0..4, Z:0..4, Q:0..4")
	X, Y, Z, Q := createVars([]int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4},
		[]int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4})
	store.AddPropagators(CreateAlldistinct(X, Y, Z, Q))
	ready := store.IsConsistent()
	ready_test(t, "Alldistinct", ready, true)
}
