package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"testing"
)

// From "Finite Domain Constraint Programming Systems",
// Christian Schulte and Mats Carlsson, 2006

// X+Y=9 and 2X+4Y=24 with primitive constraints
// X+Y=9 and X*2=CX2, Y*4=CY2,  CX2+CY2=24
func Test_simple2Primitive(t *testing.T) {
	setup()
	defer teardown()
	log("simple2Primitive: x+y=9, 2x+4y=24 with many basic constraints")
	X := core.CreateIntVarFromTo("X", store, 0, 9)
	XMin, XMax := store.GetMinMaxDomain(X)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	YMin, YMax := store.GetMinMaxDomain(Y)
	store.AddPropagator(propagator.CreateXplusYeqC(X, Y, 9)) // X+Y=9
	multiplescx2 := rangestep(XMin, XMax*2, 2)
	CX2 := core.CreateIntVarValues("CX2", store, multiplescx2)
	store.AddPropagator(propagator.CreateXmultCeqY(X, 2, CX2)) // X*2=CX2
	multiplescy4 := rangestep(YMin, YMax*4, 4)
	CY4 := core.CreateIntVarValues("CY4", store, multiplescy4)
	store.AddPropagator(propagator.CreateXmultCeqY(Y, 4, CY4)) // Y*4=CY4
	// 2*X+4*Y=24 iff CX2+CY4=24
	store.AddPropagator(propagator.CreateXplusYeqC(CX2, CY4, 24))
	ready := store.IsConsistent()
	ready_test(t, "simple2Primitive", ready, true)
	domainEquals_test(t, "simple2Primitive", X, []int{6})
	domainEquals_test(t, "simple2Primitive", Y, []int{3})
	// Note, that as all variables are ground, the propagators
	// shall no longer be registered
	equalsInt_test(t, "propagators still active",
		store.GetStat().GetActPropagators(), 0)
	propStat()
}

// X+Y=9 and 2X+4Y=24 with dedicated constraints
func Test_simple2Dedicated(t *testing.T) {
	setup()
	defer teardown()
	log("simple2Dedicated: x+y=9, 2x+4y=24 with two dedicated constraints")
	X := core.CreateIntVarFromTo("X", store, 0, 9)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	prop1 := propagator.CreateC1XplusC2YeqC3(1, X, 1, Y, 9)
	store.AddPropagator(prop1)
	prop2 := propagator.CreateC1XplusC2YeqC3(2, X, 4, Y, 24)
	store.AddPropagator(prop2)
	ready := store.IsConsistent()
	ready_test(t, "simple2Dedicated", ready, true)
	domainEquals_test(t, "simple2Dedicated", X, []int{6})
	domainEquals_test(t, "simple2Dedicated", Y, []int{3})
	// Note, that as all variables are ground, the propagators
	// shall no longer be registered
	equalsInt_test(t, "propagators still active",
		store.GetStat().GetActPropagators(), 0)
	propStat()
}
