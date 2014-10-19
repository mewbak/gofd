package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xneqc_ic_test(t *testing.T, from int, to int, c int,
	expx []int, expready bool) {
	xDom := core.CreateIvDomainFromTo(from, to)
	X := core.CreateIntVarDom("X", store, xDom)
	xgtc := CreateXneqC(X, c)
	store.AddPropagator(xgtc)
	ready := store.IsConsistent()
	ready_test(t, "XneqC_IC", ready, expready)
	if expready {
		domainEquals_test(t, "XneqC_IC", X, expx)
	}
}

func Test_XneqC_IC1(t *testing.T) {
	setup()
	defer teardown()
	log("XneqC_IC1: C=5, X:0..4")
	xneqc_ic_test(t, 0, 4, 5, []int{0, 1, 2, 3, 4}, true)
}

func Test_XneqC_IC2(t *testing.T) {
	setup()
	defer teardown()
	log("XneqC_IC2: C=6, X:6..9")
	xneqc_ic_test(t, 6, 9, 6, []int{7, 8, 9}, true)
}

func Test_XneqC_IC3(t *testing.T) {
	setup()
	defer teardown()
	log("XneqC_IC3: C=4, X:4..4")
	xneqc_ic_test(t, 4, 4, 4, []int{}, false)
}
