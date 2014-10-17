package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xNeqY_test(t *testing.T, xinit []int, yinit []int,
	expx []int, expy []int, expready bool) {
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	prop := CreateXneqY(X, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "XneqY", ready, expready)
	if expready {
		domainEquals_test(t, "XneqY", X, expx)
		domainEquals_test(t, "XneqY", Y, expy)
	}
}

// tests for X!=Y
func Test_XneqY1(t *testing.T) {
	setup()
	defer teardown()
	log("XneqY1: X!=Y, X:0..1, Y:0..1")
	xinit := []int{0, 1}
	yinit := []int{0, 1}
	expx := []int{0, 1}
	expy := []int{0, 1}
	xNeqY_test(t, xinit, yinit, expx, expy, true)
}

func Test_XneqY2(t *testing.T) {
	setup()
	defer teardown()
	log("XneqY2: X!=Y, X:0, Y:1")

	xinit := []int{0}
	yinit := []int{1}
	expx := []int{0}
	expy := []int{1}
	xNeqY_test(t, xinit, yinit, expx, expy, true)
}

func Test_XneqY3(t *testing.T) {
	setup()
	defer teardown()
	log("XneqY3: X!=Y, X:1, Y:1")

	xinit := []int{1}
	yinit := []int{1}
	expx := []int{}
	expy := []int{}
	xNeqY_test(t, xinit, yinit, expx, expy, false)
}
