package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xPlusCEqY_test(t *testing.T, xinit []int, c int, yinit []int,
	expx []int, expy []int, expready bool) {
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	prop := CreateXplusCeqY(X, c, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "XplusCeqY", ready, expready)
	if expready {
		domainEquals_test(t, "XplusCeqY", X, expx)
		domainEquals_test(t, "XplusCeqY", Y, expy)
	}
}

func Test_XplusCeqY5(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCeqY5: X+5=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 5
	expx := []int{0, 1, 2, 3, 4}
	expy := []int{5, 6, 7, 8, 9}
	xPlusCEqY_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XplusCeqY0(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCeqY0: X+0=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 0
	expx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	xPlusCEqY_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XplusCeqY10(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCeqY10: X+10=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 10
	expx := []int{}
	expy := []int{}
	xPlusCEqY_test(t, xinit, c, yinit, expx, expy, false)
}
