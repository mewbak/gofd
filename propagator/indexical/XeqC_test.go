package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xeqc_ic_test(t *testing.T, from int, to int, c int,
	expx []int, expready bool) {
	xDom := core.CreateIvDomainFromTo(from, to)
	X := core.CreateIntVarDom("X", store, xDom)
	xgtc := CreateXeqC(X, c)
	store.AddPropagator(xgtc)
	ready := store.IsConsistent()
	ready_test(t, "XeqC_IC", ready, expready)
	if expready {
		domainEquals_test(t, "XeqC_IC", X, expx)
	}
}

func Test_XeqC_IC1(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_IC1: C=5, X:0..9")
	xeqc_ic_test(t, 0, 9, 5, []int{5}, true)
}

func Test_XeqC_IC2(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_IC2: C=6, X:4..6")
	xeqc_ic_test(t, 4, 6, 6, []int{6}, true)
}

func Test_XeqC_IC3(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_IC3: C=5, X:0..4")
	xeqc_ic_test(t, 0, 4, 5, []int{}, false)
}
