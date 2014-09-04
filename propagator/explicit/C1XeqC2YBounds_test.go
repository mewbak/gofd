package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func C1XeqC2YBounds_test(t *testing.T, c1 int, xinit []int, c2 int, yinit []int,
	expx []int, expy []int, expready bool) {
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	prop := CreateC1XeqC2YBounds(c1, X, c2, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "C1XeqC2YBounds", ready, expready)

	if expready {
		domainEquals_test(t, "C1XeqC2YBounds", X, expx)
		domainEquals_test(t, "C1XeqC2YBounds", Y, expy)
	}
}

func Test_C1XeqC2YBoundsa(t *testing.T) {
	setup()
	defer teardown()
	log("C1XeqC2YBoundsa: X*3=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c1 := 3
	c2 := 1
	expx := []int{0, 1, 2, 3}
	expy := []int{0, 3, 6, 9}
	C1XeqC2YBounds_test(t, c1, xinit, c2, yinit, expx, expy, true)
}

func Test_C1XeqC2YBoundsb(t *testing.T) {
	setup()
	defer teardown()
	log("C1XeqC2YBoundsb: X*0=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c1 := 0
	c2 := 1
	expx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0}
	C1XeqC2YBounds_test(t, c1, xinit, c2, yinit, expx, expy, true)
}

func Test_C1XeqC2YBoundsc(t *testing.T) {
	setup()
	defer teardown()

	log("C1XeqC2YBoundsc: X*10=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c1 := 10
	c2 := 1
	expx := []int{0}
	expy := []int{0}
	C1XeqC2YBounds_test(t, c1, xinit, c2, yinit, expx, expy, true)
}

func Test_C1XeqC2YBoundsd(t *testing.T) {
	setup()
	defer teardown()
	log("C1XeqC2YBoundsd: X*10=Y, X:1..9, Y:1..9")

	xinit := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	c1 := 10
	c2 := 1
	expx := []int{}
	expy := []int{}
	C1XeqC2YBounds_test(t, c1, xinit, c2, yinit, expx, expy, false)
}

func Test_C1XeqC2YBoundsf(t *testing.T) {
	setup()
	defer teardown()
	log("C1XeqC2YBoundsf: X=Y, X:0..9, Y:0..9")

	xinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	yinit := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	c1 := 1
	c2 := 1
	expx := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	expy := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	C1XeqC2YBounds_test(t, c1, xinit, c2, yinit, expx, expy, true)
}
