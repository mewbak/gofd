package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xeqc_test(t *testing.T, xinit []int, c int,
	expx []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	xgtc := CreateXeqC(X, c)
	store.AddPropagator(xgtc)
	ready := store.IsConsistent()
	ready_test(t, "XeqC", ready, expready)
	if expready {
		domainEquals_test(t, "XeqC", X, expx)
	}
}

func Test_XeqC1(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC1: C=5, X:0..9")
	xeqc_test(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 5, []int{5}, true)
}

func Test_XeqC2(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC2: C=6, X:4..6")
	xeqc_test(t, []int{4, 5, 6}, 6, []int{6}, true)
}

func Test_XeqC3(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC3: C=5, X:0..4")
	xeqc_test(t, []int{0, 1, 2, 3, 4}, 5, []int{}, false)
}

func Test_XeqC_clone(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_clone")

	xinit := []int{1, 2, 3, 4, 5}
	x := core.CreateIntVarIvValues("X", store, xinit)
	c := CreateXeqC(x, 5)

	clone_test(t, store, c)
}
