package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xneqc_test(t *testing.T, xinit []int, c int,
	expx []int, expready bool) {
	X := core.CreateIntVarExValues("X", store, xinit)
	xgtc := CreateXneqC(X, c)
	store.AddPropagator(xgtc)
	ready := store.IsConsistent()
	ready_test(t, "XneqC", ready, expready)
	if expready {
		domainEquals_test(t, "XneqC", X, expx)
	}
}

func Test_XneqC5a(t *testing.T) {
	setup()
	defer teardown()
	log("XneqC5a: X!=5, X:0..9")
	xneqc_test(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 5,
		[]int{0, 1, 2, 3, 4, 6, 7, 8, 9}, true)
}

func Test_XneqC6(t *testing.T) {
	setup()
	defer teardown()
	log("XneqC6 : X!=6, X:4..6")
	xneqc_test(t, []int{4, 5, 6}, 6, []int{4, 5}, true)
}

func Test_XneqC5b(t *testing.T) {
	setup()
	defer teardown()
	log("XneqC5b: X!=5, X:0..4")
	xneqc_test(t, []int{0, 1, 2, 3, 4}, 5, []int{0, 1, 2, 3, 4}, true)
}

func Test_XneqC0(t *testing.T) {
	setup()
	defer teardown()
	log("XneqC0 : X!=0, X:0")
	xneqc_test(t, []int{0}, 0, []int{}, false)
}
