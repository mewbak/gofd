package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xPlusCEqY_IC_test(t *testing.T, xinit []int, c int, yinit []int,
	expx []int, expy []int, expready bool) {
	xDom := core.CreateIvDomainFromIntArr(xinit)
	yDom := core.CreateIvDomainFromIntArr(yinit)
	X := core.CreateIntVarDom("X", store, xDom)
	Y := core.CreateIntVarDom("Y", store, yDom)
	prop := CreateXplusCeqY(X, c, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "XplusCeqY_IC", ready, expready)
	if expready {
		domainEquals_test(t, "XplusCeqY_IC", X, expx)
		domainEquals_test(t, "XplusCeqY_IC", Y, expy)
	}
}

func Test_XplusCeqY_IC_1(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCeqY_IC_1: X+5=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 5
	expx := []int{0, 1, 2, 3, 4}
	expy := []int{5, 6, 7, 8, 9}
	xPlusCEqY_IC_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XplusCeqY_IC_2(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCeqY_IC_2: X+0=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 0
	expx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	xPlusCEqY_IC_test(t, xinit, c, yinit, expx, expy, true)
}

func Test_XplusCeqY_IC_3(t *testing.T) {
	setup()
	defer teardown()
	log("XplusCeqY_IC_3: X+10=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c := 10
	expx := []int{}
	expy := []int{}
	xPlusCEqY_IC_test(t, xinit, c, yinit, expx, expy, false)
}
