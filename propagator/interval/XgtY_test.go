package interval

import (
	"testing"
)

func xGtY_test(t *testing.T, xinit []int, yinit []int,
	expx []int, expy []int, expready bool) {
	X, Y := createXYtestVars(xinit, yinit)
	prop := CreateXgtY(X, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "XgtY", ready, expready)
	if expready {
		domainEquals_test(t, "XgtY", X, expx)
		domainEquals_test(t, "XgtY", Y, expy)
	}
}

// tests for X>Y
func Test_XgtY1(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY1: X>Y, X:0..9, Y:0..9")
	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expx := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	xGtY_test(t, xinit, yinit, expx, expy, true)
}

func Test_XgtY2(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY2: X>Y, X:0..3, Y:4..9")

	xinit := []int{0, 1, 2, 3}
	yinit := []int{4, 5, 6, 7, 8, 9}
	expx := []int{}
	expy := []int{}
	xGtY_test(t, xinit, yinit, expx, expy, false)
}

func Test_XgtY3(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY3: X>Y, X:0..4, Y:3..9")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{3, 4, 5, 6, 7, 8, 9}
	expx := []int{4}
	expy := []int{3}
	xGtY_test(t, xinit, yinit, expx, expy, true)
}

func Test_XgtY_clone(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY_clone")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{3, 4, 5, 6, 7, 8, 9}

	X, Y := createXYtestVars(xinit, yinit)
	constraint := CreateXgtY(X, Y)

	clone_test(t, store, constraint)
}
