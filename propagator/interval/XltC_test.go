package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xltc_test(t *testing.T, xinit []int, c int,
	expx []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	xgtc := CreateXltC(X, c)
	store.AddPropagator(xgtc)
	ready := store.IsConsistent()
	ready_test(t, "XltC", ready, expready)
	if expready {
		domainEquals_test(t, "XltC", X, expx)
	}
}

func xlteqc_test(t *testing.T, xinit []int, c int,
	expx []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	xgtc := CreateXlteqC(X, c)
	store.AddPropagator(xgtc)
	ready := store.IsConsistent()
	if ready != expready {
		t.Errorf("XltC: ready = %v, want %v",
			ready, expready)
	}
	XDomain := store.GetDomain(X)
	expDomain := core.CreateExDomainAdds(expx)
	if !XDomain.Equals(expDomain) {
		t.Errorf("XltC: X > %d, got X=%s, want X=%s\n",
			c, XDomain.String(), expDomain.String())
	}
}

func Test_XltC5a(t *testing.T) {
	setup()
	defer teardown()
	log("XltC5a: X<5, X:0..9")
	xltc_test(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 5,
		[]int{0, 1, 2, 3, 4}, true)
}

func Test_XltC5b(t *testing.T) {
	setup()
	defer teardown()
	log("XltC5b: X<5, X:0..6")
	xltc_test(t, []int{0, 1, 2, 3, 4, 5, 6}, 5, []int{0, 1, 2, 3, 4}, true)
}

func Test_XltC5c(t *testing.T) {
	setup()
	defer teardown()
	log("XltC5c: X<5, X:0..4")
	xltc_test(t, []int{0, 1, 2, 3, 4}, 5, []int{0, 1, 2, 3, 4}, true)
}

func Test_XlteqC5a(t *testing.T) {
	setup()
	defer teardown()
	log("XlteqC5a: X<=5, X:5..9")
	xlteqc_test(t, []int{5, 6, 7, 8, 9}, 5, []int{5}, true)
}

func Test_XlteqC5b(t *testing.T) {
	setup()
	defer teardown()
	log("XlteqC5b: X<=5, X:6..9")
	xlteqc_test(t, []int{6, 7, 8, 9}, 5, []int{}, false)
}

func Test_XltC_clone(t *testing.T) {
	setup()
	defer teardown()
	log("XltC_clone")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 5

	X := core.CreateIntVarIvValues("X", store, xinit)
	constraint := CreateXltC(X, c)

	clone_test(t, store, constraint)
}
