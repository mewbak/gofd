package interval

import (
	"testing"
)

func C1XmultC2YeqC3ZBounds_test(t *testing.T, c1 int, xinit []int,
	c2 int, yinit []int, c3 int, zinit []int,
	expx []int, expy []int, expz []int, expready bool) {
	//test_logger.SetLoggingLevel(core.LOG_DEBUG)
	X, Y, Z := createC1XplusC2YeqC3ZtestVars(xinit, yinit, zinit)
	store.AddPropagator(CreateC1XmultC2YeqC3ZBounds(c1, X, c2, Y, c3, Z))
	ready := store.IsConsistent()
	ready_test(t, "C1XmultC2YeqC3ZBounds", ready, expready)
	if expready {
		domainEquals_test(t, "C1XmultC2YeqC3ZBounds", X, expx)
		domainEquals_test(t, "C1XmultC2YeqC3ZBounds", Y, expy)
		domainEquals_test(t, "C1XmultC2YeqC3ZBounds", Z, expz)
	}
}

func Test_GC1XmultC2YeqC3ZBoundsa(t *testing.T) {
	setup()
	defer teardown()
	log("C1XmultC2YeqC3ZBoundsa: X*Y=Z, X:0..4, Y:0..4, Z:6,8,9,16")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{6, 8, 9, 16}
	c1, c2, c3 := 1, 1, 1
	expx := []int{2, 3, 4}
	expy := []int{2, 3, 4}
	expz := []int{6, 8, 9, 16}
	C1XmultC2YeqC3ZBounds_test(t, c1, xinit, c2, yinit, c3, zinit,
		expx, expy, expz, true)
}

func Test_C1XmultC2YeqC3ZBoundsb(t *testing.T) {
	setup()
	defer teardown()
	log("C1XmultC2YeqC3ZBoundsb: X+Y=Z, X:0..4, Y:0..4, Z:16")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{16}
	c1, c2, c3 := 1, 1, 1
	expx := []int{4}
	expy := []int{4}
	expz := []int{16}
	C1XmultC2YeqC3ZBounds_test(t, c1, xinit, c2, yinit, c3, zinit,
		expx, expy, expz, true)
}

func Test_C1XmultC2YeqC3ZBoundsc(t *testing.T) {
	setup()
	defer teardown()
	log("C1XmultC2YeqC3ZBoundsc: X+Y=Z, X:0, Y:0, Z:1")

	xinit := []int{0}
	yinit := []int{0}
	zinit := []int{1}
	c1, c2, c3 := 1, 1, 1
	expx := []int{}
	expy := []int{}
	expz := []int{}
	C1XmultC2YeqC3ZBounds_test(t, c1, xinit, c2, yinit, c3, zinit,
		expx, expy, expz, false)
}

func Test_C1XmultC2YeqC3ZBoundsd(t *testing.T) {
	setup()
	defer teardown()
	log("C1XmultC2YeqC3ZBoundsd: X+Y=Z, X:1..4, Y:2..4, Z:0,1")

	xinit := []int{1, 2, 3, 4}
	yinit := []int{2, 3, 4}
	zinit := []int{0, 1}
	c1, c2, c3 := 1, 1, 1
	expx := []int{}
	expy := []int{}
	expz := []int{}
	C1XmultC2YeqC3ZBounds_test(t, c1, xinit, c2, yinit, c3, zinit,
		expx, expy, expz, false)
}

func Test_C1XmultC2YeqC3ZBounds_clone(t *testing.T) {
	setup()
	defer teardown()
	log("C1XmultC2YeqC3ZBounds_clone")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{6, 8, 9, 16}
	c1, c2, c3 := 1, 1, 1

	X, Y, Z := createC1XplusC2YeqC3ZtestVars(xinit, yinit, zinit)
	c := CreateC1XmultC2YeqC3ZBounds(c1, X, c2, Y, c3, Z)

	clone_test(t, store, c)
}
