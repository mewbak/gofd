package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func XgtYIC_test(t *testing.T, xinit *core.IvDomain, yinit *core.IvDomain,
	expX *core.IvDomain, expY *core.IvDomain, expready bool) {
	X := core.CreateIntVarDom("X", store, xinit)
	Y := core.CreateIntVarDom("Y", store, yinit)
	prop := CreateXgtY(X, Y)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "XgtY_IC", ready, expready)
	if expready {
		DomainEquals_test(t, "XgtY_IC", X, expX)
		DomainEquals_test(t, "XgtY_IC", Y, expY)
	}
}

// tests for X>Y
func Test_XgtYIC1(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY_IC_1: X>Y, X:1..9, Y:1..9")

	xInit := core.CreateIvDomainFromTo(1, 9)
	yInit := core.CreateIvDomainFromTo(1, 9)

	expX := core.CreateIvDomainFromTo(2, 9)
	expY := core.CreateIvDomainFromTo(1, 8)

	XgtYIC_test(t, xInit, yInit, expX, expY, true)
}

// tests for X>Y
func Test_XgtYIC2(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY_IC_2: X>Y, X:1..20, Y:5..30")
	xInit := core.CreateIvDomainFromTo(1, 20)
	yInit := core.CreateIvDomainFromTo(5, 30)

	expX := core.CreateIvDomainFromTo(6, 20)
	expY := core.CreateIvDomainFromTo(5, 19)

	XgtYIC_test(t, xInit, yInit, expX, expY, true)
}

func Test_XgtYIC3(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY_IC_3: X>Y, X:0..3, Y:4..9")
	xInit := core.CreateIvDomainFromTo(0, 3)
	yInit := core.CreateIvDomainFromTo(4, 9)

	expX := core.CreateIvDomain()
	expY := core.CreateIvDomain()

	XgtYIC_test(t, xInit, yInit, expX, expY, false)
}

func Test_XgtYIC4(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY_IC_4: X>Y, X:1..4, Y:3..9")

	xInit := core.CreateIvDomainFromTo(1, 4)
	yInit := core.CreateIvDomainFromTo(3, 9)

	expX := core.CreateIvDomainFromTo(4, 4)
	expY := core.CreateIvDomainFromTo(3, 3)

	XgtYIC_test(t, xInit, yInit, expX, expY, true)
}

func Test_XgtYIC5(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY_IC_5: X>Y, X:1..3 6..12, Y:6..13")

	xInit := core.CreateIvDomainFromTos([][]int{{1, 3}, {6, 12}})
	yInit := core.CreateIvDomainFromTo(6, 13)

	expX := core.CreateIvDomainFromTo(7, 12)
	expY := core.CreateIvDomainFromTo(6, 11)

	XgtYIC_test(t, xInit, yInit, expX, expY, true)
}

func Test_XgtYIC6(t *testing.T) {
	setup()
	defer teardown()
	log("XgtY_IC_6: X>Y, X:1..4 6..10, 20..30, Y:5..10, 20..40")

	xInit := core.CreateIvDomainFromTos([][]int{{1, 4}, {6, 10}, {20, 30}})
	yInit := core.CreateIvDomainFromTos([][]int{{5, 10}, {20, 40}})

	expX := core.CreateIvDomainFromTos([][]int{{6, 10}, {20, 30}})
	expY := core.CreateIvDomainFromTos([][]int{{5, 10}, {20, 29}})

	XgtYIC_test(t, xInit, yInit, expX, expY, true)
}
