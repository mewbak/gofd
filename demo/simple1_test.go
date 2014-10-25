// package demo provides sample programs using gofd as finite domain
// constraint solver
package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"testing"
)

// X+Y=12 with 1*X=CX, 1*Y=CY, CX+CY=12 buildup with constants
func Test_simple1Primitive(t *testing.T) {
	setup()
	defer teardown()
	log("simple1Primitive: x+y=12 with many basic propagators")
	X := core.CreateIntVarFromTo("X", store, 0, 9)
	XMin, XMax := store.GetMinMaxDomain(X)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	YMin, YMax := store.GetMinMaxDomain(X)
	C := 12
	CX := core.CreateIntVarFromTo("CX", store, XMin, XMax*1)
	CY := core.CreateIntVarFromTo("CY", store, YMin, YMax*1)
	store.AddPropagator(propagator.CreateXmultCeqY(X, 1, CX))  // X*1=CX
	store.AddPropagator(propagator.CreateXmultCeqY(Y, 1, CY))  // Y*1=CY
	store.AddPropagator(propagator.CreateXplusYeqC(CX, CY, C)) // CX+CY=12
	ready := store.IsConsistent()
	ready_test(t, "simple1_primitve_test", ready, true)
	domainEquals_test(t, "simple1_primitive", X, rangestep(3, 9, 1))
	domainEquals_test(t, "simple1_primitive", Y, rangestep(3, 9, 1))
	propStat()
}

// X+Y=12 as one dedicated contraint with constants
func Test_simple1Dedicated(t *testing.T) {
	setup()
	defer teardown()
	log("simple1Dedicated: x+y=12 with one dedicated propagator")
	X := core.CreateIntVarFromTo("X", store, 0, 9)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	prop := propagator.CreateC1XplusC2YeqC3(1, X, 1, Y, 12)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "simple1_dedicated_test", ready, true)
	domainEquals_test(t, "simple1_dedicated", X, rangestep(3, 9, 1))
	domainEquals_test(t, "simple1_dedicated", Y, rangestep(3, 9, 1))
	propStat()
}
