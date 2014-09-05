package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func c1XPlusC2YeqC3_test(t *testing.T, c1 int, xinit []int,
	c2 int, yinit []int, c3 int,
	expx []int, expy []int, expready bool) {
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	prop := CreateC1XplusC2YeqC3(c1, X, c2, Y, c3)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "C1XplusC2YeqC3", ready, expready)
	if expready {
		domainEquals_test(t, "C1XplusC2YeqC3", X, expx)
		domainEquals_test(t, "C1XplusC2YeqC3", Y, expy)
	}
}

func Test_C1XplusC2YeqC3a(t *testing.T) {
	setup()
	defer teardown()
	log("C1XplusC2YeqC3a: 1*X+2*Y=12, X:0..9, Y:0..9")
	c1 := 1
	c2 := 2
	c3 := 12
	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expx := []int{6, 4, 8, 0, 2}
	expy := []int{4, 2, 5, 6, 3}
	c1XPlusC2YeqC3_test(t, c1, xinit, c2, yinit, c3, expx, expy, true)
}

func Test_C1XplusC2YeqC3b(t *testing.T) {
	setup()
	defer teardown()
	log("C1XplusC2YeqC3b: 1*X+2*Y=0, X:0..9, Y:0..9")

	c1 := 1
	c2 := 2
	c3 := 0
	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expx := []int{0}
	expy := []int{0}
	c1XPlusC2YeqC3_test(t, c1, xinit, c2, yinit, c3, expx, expy, true)
}

func Test_C1XplusC2YeqC3c(t *testing.T) {
	setup()
	defer teardown()
	log("C1XplusC2YeqC3c: 2*X+4*Y=3, X:0..9, Y:0..9")

	c1 := 2
	c2 := 4
	c3 := 3
	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expx := []int{}
	expy := []int{}
	c1XPlusC2YeqC3_test(t, c1, xinit, c2, yinit, c3, expx, expy, false)
}
