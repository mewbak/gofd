package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"testing"
)

func xplusyeqz_reification_test(t *testing.T, fromX int, toX int, fromY int, toY int,
	fromZ int, toZ int, fromB int, toB int,
	expx []int, expy []int, expz []int, expB []int, expready bool) {

	xDom := core.CreateIvDomainFromTo(fromX, toX)
	yDom := core.CreateIvDomainFromTo(fromY, toY)
	zDom := core.CreateIvDomainFromTo(fromZ, toZ)

	bDom := core.CreateIvDomainFromTo(fromB, toB)
	X := core.CreateIntVarDom("X", store, xDom)
	Y := core.CreateIntVarDom("Y", store, yDom)
	Z := core.CreateIntVarDom("Z", store, zDom)
	B := core.CreateIntVarDom("B", store, bDom)

	xplusyeqz := indexical.CreateXplusYeqZ(X, Y, Z)

	ricC := CreateReifiedConstraint(xplusyeqz, B)

	store.AddPropagator(ricC)
	ready := store.IsConsistent()
	ready_test(t, "XplusYeqZ_reification", ready, expready)
	if expready {
		domainEquals_test(t, "XplusYeqZ_reification", X, expx)
		domainEquals_test(t, "XplusYeqZ_reification", Y, expy)
		domainEquals_test(t, "XplusYeqZ_reification", Z, expz)
		domainEquals_test(t, "XplusYeqZ_reification", B, expB)
	}
}

//delayed
func Test_XplusYeqZ_reification1(t *testing.T) {
	setup()
	defer teardown()
	//core.GetLogger().SetLoggingLevel(core.LOG_DEBUG)
	log("XplusYeqZ_reification1: X+Y=Z<=>B, X:0..5, Y:0..5, Z:9, B:0,1")
	xplusyeqz_reification_test(t, 0, 5, 0, 5, 9, 9, 0, 1, []int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5},
		[]int{9}, []int{0, 1}, true)
}

//C disentailed bzw. !C entailed
func Test_XplusYeqZ_reification2(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZ_reification2: X+Y=Z<=>B, X:5, Y:5, Z=11, B:0,1")
	xplusyeqz_reification_test(t, 5, 5, 5, 5, 11, 11, 0, 1, []int{5}, []int{5}, []int{11},
		[]int{0}, true)
}

//C entailed
func Test_XplusYeqZ_reification3(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZ_reification3: X+Y=Z<=>B, X:0..5, Y:0..5, Z:10, B:0,1")
	xplusyeqz_reification_test(t, 0, 5, 0, 5, 10, 10, 0, 1, []int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5}, []int{10}, []int{0, 1}, true)
}

func Test_XplusYeqZ_reification4(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZ_reification4: X+Y=Z<=>B, X:0..5, Y:0..5, Z:10, B:0")
	xplusyeqz_reification_test(t, 0, 5, 0, 5, 10, 10, 0, 0, []int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5}, []int{10}, []int{0}, true)
}

//B=1
func Test_XplusYeqZ_reification5(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYeqZ_reification5: X+Y=Z<=>B, X:0..5, Y:0..5, Z:10, B:1")
	xplusyeqz_reification_test(t, 0, 5, 0, 5, 10, 10, 1, 1, []int{5}, []int{5}, []int{10}, []int{1}, true)
}
