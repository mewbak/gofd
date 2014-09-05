package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xGtEqY_test(t *testing.T, xinit []int, yinit []int,
	expx []int, expy []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	Y := core.CreateIntVarIvValues("Y", store, yinit)
	prop := CreateXgteqY(X, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	if ready != expready {
		t.Errorf("XgteqY_test: XgteqY_test-result = %v, want %v",
			ready, !ready)
	}
	if expready {
		dx := core.CreateExDomain()
		dx.Adds(expx)
		dy := core.CreateExDomain()
		dy.Adds(expy)
		msg := "XgteqY_test-DomainCheck: Domain calculated = %v, want %v"
		if !store.GetDomain(X).Equals(dx) {
			t.Errorf(msg, store.GetDomain(X), dx)
		}
		if !store.GetDomain(Y).Equals(dy) {
			t.Errorf(msg, store.GetDomain(Y), dy)
		}
	}
}

func xGtYPlusC_test(t *testing.T, xinit []int, yinit []int, C int,
	expx []int, expy []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	Y := core.CreateIntVarIvValues("Y", store, yinit)
	prop := CreateXgtYplusC(X, Y, C)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	if ready != expready {
		t.Errorf("XgtYplusC_test: XgtYplusC_test-result = %v, want %v",
			ready, !ready)
	}
	if expready {
		dx := core.CreateExDomain()
		dx.Adds(expx)
		dy := core.CreateExDomain()
		dy.Adds(expy)
		msg := "XgtYplusC_test-DomainCheck: Domain calculated = %v, want %v"
		if !store.GetDomain(X).Equals(dx) {
			t.Errorf(msg, store.GetDomain(X), dx)
		}
		if !store.GetDomain(Y).Equals(dy) {
			t.Errorf(msg, store.GetDomain(Y), dy)
		}
	}
}

// tests for X>=Y
func Test_XgteqY1(t *testing.T) {
	setup()
	defer teardown()
	log("XgteqY1: X>=Y, X:5..9, Y:0..9")

	xinit := []int{5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expx := []int{5, 6, 7, 8, 9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	xGtEqY_test(t, xinit, yinit, expx, expy, true)
}

func Test_XgteqY2(t *testing.T) {
	setup()
	defer teardown()
	log("XgteqY2: X>=Y, X:0..3, Y:4..9")

	xinit := []int{0, 1, 2, 3}
	yinit := []int{4, 5, 6, 7, 8, 9}
	expx := []int{}
	expy := []int{}
	xGtEqY_test(t, xinit, yinit, expx, expy, false)
}

func Test_XgteqY3(t *testing.T) {
	setup()
	defer teardown()
	log("XgteqY3: X>=Y, X:0..3, Y:3..9")

	xinit := []int{0, 1, 2, 3}
	yinit := []int{3, 4, 5, 6, 7, 8, 9}
	expx := []int{3}
	expy := []int{3}
	xGtEqY_test(t, xinit, yinit, expx, expy, true)
}

// tests for X>Y+C
func Test_XgtYplusC1(t *testing.T) {
	setup()
	defer teardown()
	log("XgtYplusC1: X>Y+5, X:5..9, Y:0..9")

	xinit := []int{5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	C := 5
	expx := []int{6, 7, 8, 9}
	expy := []int{0, 1, 2, 3}
	xGtYPlusC_test(t, xinit, yinit, C, expx, expy, true)
}

func Test_XgtYplusC2(t *testing.T) {
	setup()
	defer teardown()
	log("XgtYplusC2: X>Y+5, X:0..3, Y:4..9")

	xinit := []int{0, 1, 2, 3}
	yinit := []int{4, 5, 6, 7, 8, 9}
	C := 5
	expx := []int{}
	expy := []int{}
	xGtYPlusC_test(t, xinit, yinit, C, expx, expy, false)
}

func Test_XgtYplusC3(t *testing.T) {
	setup()
	defer teardown()
	log("XgtYplusC3: X>Y+5, X:0..6, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	C := 5
	expx := []int{6}
	expy := []int{0}
	xGtYPlusC_test(t, xinit, yinit, C, expx, expy, true)
}
