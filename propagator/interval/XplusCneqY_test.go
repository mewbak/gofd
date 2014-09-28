package interval

import (
	"testing"
)

func ivxpluscneqy_test(t *testing.T, xinit []int, c int, yinit []int,
	expx []int, expy []int, expready bool) {
	X, Y := createXYtestVars(xinit, yinit)
	prop := CreateXplusCneqY(X, c, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "XplusCneqY_intervals", ready, expready)
	if expready {
		domainEquals_test(t, "XplusCneqY_intervals", X, expx)
		domainEquals_test(t, "XplusCneqY_intervals", Y, expy)
	}
}

func Test_XplusCneqY1(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCneqY1_intervals: X+3!=Y, X:0..9, Y:0..9")
	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 3
	expx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	ivxpluscneqy_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XplusCneqY2(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCneqY2_intervals: X+3!=Y, X:6, Y:0..9")

	xinit := []int{6}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 3
	expx := []int{6}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	ivxpluscneqy_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XplusCneqY3(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCneqY3_intervals: X!=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 0
	expx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	ivxpluscneqy_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XplusCneqY4(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCneqY4_intervals: X!=Y, X:9, Y:0..9")

	xinit := []int{9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 0
	expx := []int{9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	ivxpluscneqy_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XplusCneqY_clone(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCneqY_clone")

	xinit := []int{9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 0

	X, Y := createXYtestVars(xinit, yinit)
	constraint := CreateXplusCneqY(X, c, Y)

	clone_test(t, store, constraint)
}
