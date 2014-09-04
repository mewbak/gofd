package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xgtc_test(t *testing.T, xinit []int, c int,
	expx []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	xgtc := CreateXgtC(X, c)
	store.AddPropagator(xgtc)
	ready := store.IsConsistent()
	ready_test(t, "XgtC", ready, expready)
	if expready {
		domainEquals_test(t, "XgtC", X, expx)
	}
}

func xgteqc_test(t *testing.T, xinit []int, c int,
	expx []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	xgtc := CreateXgteqC(X, c)
	store.AddPropagator(xgtc)
	ready := store.IsConsistent()
	if ready != expready {
		t.Errorf("XgteqC: ready = %v, want %v",
			ready, expready)
	}
	if expready {
		XDomain := store.GetDomain(X)
		expDomain := core.CreateExDomainAdds(expx)
		if !XDomain.Equals(expDomain) {
			t.Errorf("XgteqC: X > %d, got X=%s, want X=%s\n",
				c, XDomain.String(), expDomain.String())
		}
	}
}

func Test_XgtC5a(t *testing.T) {
	setup()
	defer teardown()
	log("XgtC5a: X>5, X:0..9")
	xgtc_test(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 5,
		[]int{6, 7, 8, 9}, true)
}

func Test_XgtC5b(t *testing.T) {
	setup()
	defer teardown()
	log("XgtC5b: X>5, X:0..6")
	xgtc_test(t, []int{0, 1, 2, 3, 4, 5, 6}, 5, []int{6}, true)
}

func Test_XgtC5c(t *testing.T) {
	setup()
	defer teardown()
	log("XgtC5c: X>5, X:0..4")
	xgtc_test(t, []int{0, 1, 2, 3, 4}, 5, []int{}, false)
}

func Test_XgteqC4(t *testing.T) {
	setup()
	defer teardown()
	log("XgteqC4: X>=4, X:0..4")
	xgteqc_test(t, []int{0, 1, 2, 3, 4}, 4, []int{4}, true)
}

func Test_XgteqC5(t *testing.T) {
	setup()
	defer teardown()
	log("XgteqC5: X>=5, X:0..4")
	xgteqc_test(t, []int{0, 1, 2, 3, 4}, 5, []int{}, false)
}

// fails with deadlock on copied store although
// * almost nothing happens
// * all channels are closed
func Test_Xgteq_storeclone(t *testing.T) {
	setup()
	defer teardown()
	log("XgteqC: copy store ")

	nstore := core.CreateStore()
	xinit := []int{0, 1, 2, 3, 4}
	c := 0
	X := core.CreateIntVarIvValues("X", nstore, xinit)
	xgtc := CreateXgtC(X, c)
	nstore.AddPropagator(xgtc)
	var nnstore *core.Store = nil
	for i := 0; i < 10; i++ {
		nnstore = nstore.Clone(nil)
	}
	nnstore.IsConsistent()
}
