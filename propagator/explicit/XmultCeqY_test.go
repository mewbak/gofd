package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xMultCEqY_test(t *testing.T, xinit []int, c int, yinit []int,
	expx []int, expy []int, expready bool) {
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	prop := CreateXmultCeqY(X, c, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "XmultCeqY", ready, expready)
	if expready {
		domainEquals_test(t, "XmultCeqY", X, expx)
		domainEquals_test(t, "XmultCeqY", Y, expy)
	}
}

func Test_XmultCeqY3(t *testing.T) {
	setup()
	defer teardown()
	log("XmultCeqY3: X*3=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 3
	expx := []int{0, 1, 2, 3}
	expy := []int{0, 3, 6, 9}
	xMultCEqY_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XmultCeqY0(t *testing.T) {
	setup()
	defer teardown()
	log("XmultCeqY0: X*0=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 0
	expx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0}
	xMultCEqY_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XmultCeqY10a(t *testing.T) {
	setup()
	defer teardown()
	log("XmultCeqY10a: X*10=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 10
	expx := []int{0}
	expy := []int{0}
	xMultCEqY_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XmultCeqY10b(t *testing.T) {
	setup()
	defer teardown()
	log("XmultCeqY10b: X*10=Y, X:1..9, Y:1..9")

	xinit := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 10
	expx := []int{}
	expy := []int{}
	xMultCEqY_test(t, xinit, c, yinit, expx, expy, false)
}

func Test_XmultCeqY1(t *testing.T) {
	setup()
	defer teardown()
	log("XmultCeqY1: X*10=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 1
	expx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	xMultCEqY_test(t, xinit, c, yinit, expx, expy, true)
}
