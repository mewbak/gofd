package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"testing"
)

func xeqc_reification_test(t *testing.T, fromX int, toX int, c int,
	fromB int, toB int,
	expx []int, expB []int, expready bool) {
	xDom := core.CreateIvDomainFromTo(fromX, toX)
	bDom := core.CreateIvDomainFromTo(fromB, toB)
	X := core.CreateIntVarDom("X", store, xDom)
	B := core.CreateIntVarDom("B", store, bDom)
	xeqc := indexical.CreateXeqC(X, c)

	ricC := CreateReifiedConstraint(xeqc, B)
	store.AddPropagator(ricC)
	ready := store.IsConsistent()
	ready_test(t, "XeqC_reification", ready, expready)
	if expready {
		domainEquals_test(t, "XeqC_reification", X, expx)
		domainEquals_test(t, "XeqC_reification", B, expB)
	}
}

//delayed
func Test_XeqC_reification1(t *testing.T) {
	setup()
	defer teardown()
	//core.GetLogger().SetLoggingLevel(core.LOG_DEBUG)
	log("XeqC_reification1: X=C<=>B, X:0..5, C=5, B:0,1")
	xeqc_reification_test(t, 0, 5, 5, 0, 1, []int{0, 1, 2, 3, 4, 5}, []int{0, 1}, true)
}

//delayed
func Test_XeqC_reification2(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification2: X=C<=>B, X:4..6, C=4, B:0,1")
	xeqc_reification_test(t, 4, 6, 4, 0, 1, []int{4, 5, 6}, []int{0, 1}, true)
}

//C entailed
func Test_XeqC_reification3(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification3: X=C<=>B, X:4..4, C=4, B:0,1")
	xeqc_reification_test(t, 4, 4, 4, 0, 1, []int{4}, []int{1}, true)
}

func Test_XeqC_reification4(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification4: X=C<=>B, X:4..4, C=4, B:0")
	xeqc_reification_test(t, 4, 4, 4, 0, 0, []int{4}, []int{}, false)
}

//!C entailed
func Test_XeqC_reification5(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification5: X=C<=>B, X:0..4, C=5, B:0,1")
	xeqc_reification_test(t, 0, 4, 5, 0, 1, []int{0, 1, 2, 3, 4}, []int{0}, true)
}

func Test_XeqC_reification6(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification6: X=C<=>B, X:0..4, C=5, B:1")
	xeqc_reification_test(t, 0, 4, 5, 1, 1, []int{0, 1, 2, 3, 4}, []int{}, false)
}

//B=1
func Test_XeqC_reification7(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification7: X=C<=>B, X:4..4, C=4, B:1")
	xeqc_reification_test(t, 4, 4, 4, 1, 1, []int{4}, []int{1}, true)
}

func Test_XeqC_reification8(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification8: X=C<=>B, X:4..4, C=4, B:0")
	xeqc_reification_test(t, 4, 4, 4, 0, 0, []int{4}, []int{}, false)
}

//B=0
func Test_XeqC_reification9(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification9: X=C<=>B, X:0..4, C=5, B:0,1")
	xeqc_reification_test(t, 0, 4, 5, 0, 1, []int{0, 1, 2, 3, 4}, []int{0}, true)
}

func Test_XeqC_reification10(t *testing.T) {
	setup()
	defer teardown()
	log("XeqC_reification10: X=C<=>B, X:0..4, C=5, B:1")
	xeqc_reification_test(t, 0, 4, 5, 1, 1, []int{0, 1, 2, 3, 4}, []int{}, false)
}
