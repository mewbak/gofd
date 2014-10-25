// package demo provides sample programs using gofd as finite domain
// constraint solver
package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"testing"
)

// X+Y=12 with 1*X=CX, 1*Y=CY, CX+CY=12 buildup with constants
func Test_simple1_primitive(t *testing.T) {
	setup()
	defer teardown()
	log("simple1_primitive: x+y=12 with many basic propagators")
	X := core.CreateIntVarFromTo("X", store, 0, 9)
	XMin, XMax := store.GetMinMaxDomain(X)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	YMin, YMax := store.GetMinMaxDomain(X)
	C := 12
	CX := core.CreateIntVarFromTo("CX", store, XMin, XMax*1)
	CY := core.CreateIntVarFromTo("CY", store, YMin, YMax*1)
	store.AddPropagators(propagator.CreateXmultCeqY(X, 1, CX))  // X*1=CX
	store.AddPropagators(propagator.CreateXmultCeqY(Y, 1, CY))  // Y*1=CY
	store.AddPropagators(propagator.CreateXplusYeqC(CX, CY, C)) // CX+CY=12
	ready := store.IsConsistent()
	ready_test(t, "simple1_primitve_test", ready, true)
	domainEquals_test(t, "simple1_primitive", X, []int{3, 4, 5, 6, 7, 8, 9})
	domainEquals_test(t, "simple1_primitive", Y, []int{3, 4, 5, 6, 7, 8, 9})
	propStat()
}

// X+Y=12 as one dedicated contraint with constants
func Test_simple1_dedicated(t *testing.T) {
	setup()
	defer teardown()
	log("simple1_dedicated: x+y=12 with one dedicated propagator")
	X := core.CreateIntVarFromTo("X", store, 0, 9)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	prop := propagator.CreateC1XplusC2YeqC3(1, X, 1, Y, 12)
	store.AddPropagator(prop)
	ready := store.IsConsistent()
	ready_test(t, "simple1_dedicated_test", ready, true)
	domainEquals_test(t, "simple1_dedicated", X, []int{3, 4, 5, 6, 7, 8, 9})
	domainEquals_test(t, "simple1_dedicated", Y, []int{3, 4, 5, 6, 7, 8, 9})
	propStat()
}
