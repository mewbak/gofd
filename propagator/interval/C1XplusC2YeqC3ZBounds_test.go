package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func c1XplusC2YeqC3ZBounds_test(t *testing.T, c1 int, xinit []int, c2 int,
	yinit []int, c3 int, zinit []int,
	expx []int, expy []int, expz []int, expready bool) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	Y := core.CreateIntVarIvValues("Y", store, yinit)
	Z := core.CreateIntVarIvValues("Z", store, zinit)
	store.AddPropagator(CreateC1XplusC2YeqC3ZBounds(c1, X, c2, Y, c3, Z))
	ready := store.IsConsistent()
	ready_test(t, "C1XplusC2YeqC3ZBounds", ready, expready)
	if expready {
		domainEquals_test(t, "C1XplusC2YeqC3ZBounds", X, expx)
		domainEquals_test(t, "C1XplusC2YeqC3ZBounds", Y, expy)
		domainEquals_test(t, "C1XplusC2YeqC3ZBounds", Z, expz)
	}
}

func Test_GC1XplusC2YeqC3Za(t *testing.T) {
	setup()
	defer teardown()
	//core.GetLogger().SetLoggingLevel(core.LOG_DEBUG)
	log("C1XplusC2YeqC3ZBoundsa: X+Y=Z, X:0..4, Y:0..4, Z:6,8,9")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{6, 8, 9}
	c1, c2, c3 := 1, 1, 1
	expx := []int{2, 3, 4}
	expy := []int{2, 3, 4}
	expz := []int{6, 8}
	c1XplusC2YeqC3ZBounds_test(t, c1, xinit, c2, yinit, c3, zinit,
		expx, expy, expz, true)
}

func Test_GC1XplusC2YeqC3Zb(t *testing.T) {
	setup()
	defer teardown()
	log("C1XplusC2YeqC3ZBoundsb: X+Y=Z, X:0..4, Y:0..4, Z:6,8,9")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{1}
	c1, c2, c3 := 1, 1, 1
	expx := []int{0, 1}
	expy := []int{0, 1}
	expz := []int{1}
	c1XplusC2YeqC3ZBounds_test(t, c1, xinit, c2, yinit, c3, zinit,
		expx, expy, expz, true)
}

func Test_GC1XplusC2YeqC3Zc(t *testing.T) {
	setup()
	defer teardown()
	log("C1XplusC2YeqC3ZBoundsc: X+Y=Z, X:0, Y:0, Z:1")

	xinit := []int{0}
	yinit := []int{0}
	zinit := []int{1}
	c1, c2, c3 := 1, 1, 1
	expx := []int{}
	expy := []int{}
	expz := []int{}
	c1XplusC2YeqC3ZBounds_test(t, c1, xinit, c2, yinit, c3, zinit,
		expx, expy, expz, false)
}

func Test_GC1XplusC2YeqC3Zd(t *testing.T) {
	setup()
	defer teardown()
	log("C1XplusC2YeqC3ZBoundsd: X+Y=Z, X:0..4, Y:0..4, Z:9,10")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{9, 10}
	c1, c2, c3 := 1, 1, 1
	expx := []int{}
	expy := []int{}
	expz := []int{}
	c1XplusC2YeqC3ZBounds_test(t, c1, xinit, c2, yinit, c3, zinit,
		expx, expy, expz, false)
}

func Test_GC1XplusC2YeqC3Ze(t *testing.T) {
	setup()
	defer teardown()
	log("C1XplusC2YeqC3ZBoundse: X+X=Z, X:0..4, Z:8,9")

	xinit := []int{0, 1, 2, 3, 4}
	zinit := []int{8, 9}
	c1, c2, c3 := 1, 1, 1
	expx := []int{4}
	expz := []int{8}
	c1XplusC2YeqC3ZBounds_test(t, c1, xinit, c2, xinit, c3, zinit,
		expx, expx, expz, true)
}

func Test_GC1XplusC2YeqC3Zf(t *testing.T) {
	setup()
	defer teardown()
	log("C1XplusC2YeqC3ZBoundsf: X+X=Z, X:0..4, Z:9")

	xinit := []int{0, 1, 2, 3, 4}
	zinit := []int{9}
	c1, c2, c3 := 1, 1, 1
	expx := []int{}
	expz := []int{}
	c1XplusC2YeqC3ZBounds_test(t, c1, xinit, c2, xinit, c3, zinit,
		expx, expx, expz, false)
}
