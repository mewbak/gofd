package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func xPlusYeqZ_test(t *testing.T, xinit []int, yinit []int, zinit []int,
	expx []int, expy []int, expz []int, expready bool) {
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	Z := core.CreateIntVarExValues("Z", store, zinit)
	store.AddPropagator(CreateXplusYeqZ(X, Y, Z))
	ready := store.IsConsistent()
	ready_test(t, "XplusYeqZ", ready, expready)
	if expready {
		domainEquals_test(t, "XplusYeqZ", X, expx)
		domainEquals_test(t, "XplusYeqZ", Y, expy)
		domainEquals_test(t, "XplusYeqZ", Z, expz)
	}
}

func Test_XplusYeqZa(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZa: X+Y=Z, X:0..4, Y:0..4, Z:6,8,9")
	//core.GetLogger().SetLoggingLevel(core.LOG_DEBUG)

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{6, 8, 9}

	expx := []int{2, 3, 4}
	expy := []int{2, 3, 4}
	expz := []int{6, 8}

	xPlusYeqZ_test(t, xinit, yinit, zinit, expx, expy, expz, true)
}

func Test_XplusYeqZb(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZb: X+Y=Z, X:0..4, Y:0..4, Z:6,8,9")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{1}

	expx := []int{0, 1}
	expy := []int{0, 1}
	expz := []int{1}

	xPlusYeqZ_test(t, xinit, yinit, zinit, expx, expy, expz, true)
}

func Test_XplusYeqZc(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZc: X+Y=Z, X:0, Y:0, Z:1")

	xinit := []int{0}
	yinit := []int{0}
	zinit := []int{1}

	expx := []int{}
	expy := []int{}
	expz := []int{}

	xPlusYeqZ_test(t, xinit, yinit, zinit, expx, expy, expz, false)
}

func Test_XplusYeqZd(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZd: X+Y=Z, X:0..4, Y:0..4, Z:9,10")

	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{9, 10}

	expx := []int{}
	expy := []int{}
	expz := []int{}

	xPlusYeqZ_test(t, xinit, yinit, zinit, expx, expy, expz, false)
}

func Test_XplusYeqZe(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZe: X+X=Z, X:0..4, Z:8,9")

	xinit := []int{0, 1, 2, 3, 4}
	zinit := []int{8, 9}

	expx := []int{4}
	expz := []int{8}

	xPlusYeqZ_test(t, xinit, xinit, zinit, expx, expx, expz, true)
}

func Test_XplusYeqZf(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZf: X+X=Z, X:0..4, Z:9")

	xinit := []int{0, 1, 2, 3, 4}
	zinit := []int{9}

	expx := []int{}
	expz := []int{}

	xPlusYeqZ_test(t, xinit, xinit, zinit, expx, expx, expz, false)
}
