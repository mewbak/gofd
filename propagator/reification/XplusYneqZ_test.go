package reification

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"testing"
)

func XplusYneqZ_test(t *testing.T, fromX int, toX int, fromY int, toY int,
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

	xPlusYNeqZ := indexical.CreateXplusYneqZ(X, Y, Z)

	ricC := CreateReifiedConstraint(xPlusYNeqZ, B)

	store.AddPropagator(ricC)
	ready := store.IsConsistent()

	ready_test(t, "XplusYneqZreification", ready, expready)
	if expready {
		domainEquals_test(t, "XplusYneqZreification", X, expx)
		domainEquals_test(t, "XplusYneqZreification", Y, expy)
		domainEquals_test(t, "XplusYneqZreification", Z, expz)
		domainEquals_test(t, "XplusYeqZreification", B, expB)
	}
}

//delayed
func Test_XplusYneqZreification1(t *testing.T) {
	setup()
	defer teardown()
	//core.GetLogger().SetLoggingLevel(core.LOG_DEBUG)
	log("XplusYneqZreification1: X+Y!=Z<=>B, X:0..5, Y:0..5, Z:9, B:0,1")
	XplusYneqZ_test(t, 0, 5, 0, 5, 9, 9, 0, 1, []int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5},
		[]int{9}, []int{0, 1}, true)
}

//C disentailed, !C entailed
func Test_XplusYneqZreification2(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYneqZreification2: X+Y!=Z<=>B, X:5, Y:5, Z=10, B:0,1")
	XplusYneqZ_test(t, 5, 5, 5, 5, 10, 10, 0, 1, []int{5}, []int{5}, []int{10},
		[]int{0}, true)
}

//C entailed
func Test_XplusYneqZreification3(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYneqZreification3: X+Y!=Z<=>B, X:0..5, Y:0..5, Z:10, B:0,1")
	XplusYneqZ_test(t, 0, 5, 0, 5, 10, 10, 0, 1, []int{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3, 4, 5}, []int{10}, []int{0, 1}, true)
}

func Test_XplusYneqZreification4(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYneqZreification4: X+Y!=Z<=>B, X:0, Y:0, Z:10, B:0,1")
	XplusYneqZ_test(t, 0, 0, 0, 0, 10, 10, 0, 1, []int{0}, []int{0}, []int{10}, []int{1}, true)
}

func Test_XplusYneqZreification5(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYneqZ_reification5: X+Y!=Z<=>B, X:0, Y:0, Z:10, B:0")
	XplusYneqZ_test(t, 0, 0, 0, 0, 10, 10, 0, 0, []int{}, []int{}, []int{}, []int{}, false)
}

func Test_XplusYneqZreification6(t *testing.T) {
	setup()
	defer teardown()
	log("XplusYneqZreification6: X+Y!=Z<=>B, X:0, Y:0, Z:10, B:1")
	XplusYneqZ_test(t, 0, 0, 0, 0, 10, 10, 1, 1, []int{0}, []int{0}, []int{10}, []int{1}, true)
}
