package interval

import (
	"testing"
)

func xMultYeqZ_test(t *testing.T, xinit []int, yinit []int, zinit []int,
	expx []int, expy []int, expz []int, expready bool) {
	X, Y, Z := createXYZtestVars(xinit, yinit, zinit)
	store.AddPropagator(CreateXmultYeqZ(X, Y, Z))
	ready := store.IsConsistent()
	ready_test(t, "XmultYeqZ", ready, expready)
	if expready {
		domainEquals_test(t, "XmultYeqZ", X, expx)
		domainEquals_test(t, "XmultYeqZ", Y, expy)
		domainEquals_test(t, "XmultYeqZ", Z, expz)
	}
}

func Test_XmultYeqZa(t *testing.T) {
	setup()
	defer teardown()
	log("XmultYeqZa: X*Y=Z, X:0..4, Y:0..4, Z:6,8,9,16")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{6, 8, 9, 16}

	expx := []int{2, 3, 4}
	expy := []int{2, 3, 4}
	expz := []int{6, 8, 9, 16}

	xMultYeqZ_test(t, xinit, yinit, zinit, expx, expy, expz, true)
}

func Test_XmultYeqZb(t *testing.T) {
	setup()
	defer teardown()
	log("XmultYeqZb: X+Y=Z, X:0..4, Y:0..4, Z:16")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{16}

	expx := []int{4}
	expy := []int{4}
	expz := []int{16}

	xMultYeqZ_test(t, xinit, yinit, zinit, expx, expy, expz, true)
}

func Test_XmultYeqZc(t *testing.T) {
	setup()
	defer teardown()
	log("XmultYeqZc: X+Y=Z, X:0, Y:0, Z:1")

	xinit := []int{0}
	yinit := []int{0}
	zinit := []int{1}

	expx := []int{}
	expy := []int{}
	expz := []int{}

	xMultYeqZ_test(t, xinit, yinit, zinit, expx, expy, expz, false)
}

func Test_XmultYeqZd(t *testing.T) {
	setup()
	defer teardown()
	log("XmultYeqZd: X+Y=Z, X:1..4, Y:2..4, Z:0,1")

	xinit := []int{1, 2, 3, 4}
	yinit := []int{2, 3, 4}
	zinit := []int{0, 1}

	expx := []int{}
	expy := []int{}
	expz := []int{}

	xMultYeqZ_test(t, xinit, yinit, zinit, expx, expy, expz, false)
}

func Test_XmultYeqZ_clone(t *testing.T) {
	setup()
	defer teardown()
	log("XmultYeqZ_clone")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{6, 8, 9, 16}

	X, Y, Z := createXYZtestVars(xinit, yinit, zinit)
	c := CreateXmultYeqZ(X, Y, Z)

	clone_test(t, store, c)
}
