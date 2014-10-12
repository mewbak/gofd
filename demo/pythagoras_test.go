package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator"
	"fmt"
	"testing"
)

// Pythagorean triple, problem description from
// 			http://mozart.github.io/mozart-v1/doc-1.4.0/fdt/node32.html
// How many triples (A,B,C) exists, such that
//    A*A+B*B=C*C and A<=B<=C<=N, N: is a natural number
// For more test-examples: see calculator on website
//			http://www.arndt-bruenner.de/mathe/scripts/pythagotripel.htm

// C is the parameter of the equation
func testPythForSpecificC(t *testing.T, Cval int, N int, noSols int) {
	msg := "Pythagorean triple: A**2 + B**2 = %2d**2, A,B <= %2d"
	log(fmt.Sprintf(msg, Cval, N))
	A := core.CreateIntVarFromTo("A", store, 1, N)
	B := core.CreateIntVarFromTo("B", store, 1, N)
	C := core.CreateIntVarFromTo("C", store, Cval, Cval)
	AA := core.CreateIntVarFromTo("AA", store, 1, N*N)
	BB := core.CreateIntVarFromTo("BB", store, 1, N*N)
	CC := core.CreateIntVarFromTo("CC", store, N*N, N*N)
	store.AddPropagator(propagator.CreateXmultYeqZ(A, A, AA))
	store.AddPropagator(propagator.CreateXmultYeqZ(B, B, BB))
	store.AddPropagator(propagator.CreateXmultYeqZ(C, C, CC))
	store.AddPropagator(propagator.CreateXplusYeqZ(AA, BB, CC))
	store.AddPropagator(propagator.CreateXgteqY(B, A))
	store.AddPropagator(propagator.CreateXgteqY(C, B))
	query := labeling.CreateSearchAllQuery()
	result := labeling.Labeling(store, query,
		labeling.SmallestDomainFirst, labeling.InDomainMin)
	ready_test(t, "pythagoras", result, noSols > 0)
	if len(query.GetResultSet()) != noSols {
		t.Errorf("pythagorean triple, number of solutions = %d, want %d",
			len(query.GetResultSet()), noSols)
		for i, result := range query.GetResultSet() {
			t.Errorf("  %d: %d**2 + %d**2 = %d**2\n",
				i, result[A], result[B], result[C])
		}
	}
	if logger.GetLoggingLevel() <= core.LOG_NONE {
		searchStat(query.GetSearchStatistics())
	}
}

func Test_pythagorasCeq4(t *testing.T) {
	setup()
	defer teardown()
	testPythForSpecificC(t, 4, 20, 0)
}

func Test_pythagorasCeq5a(t *testing.T) {
	setup()
	defer teardown()
	testPythForSpecificC(t, 5, 5, 1)
}

func Test_pythagorasCeq5b(t *testing.T) {
	setup()
	defer teardown()
	testPythForSpecificC(t, 5, 30, 1) //slower
}

func Test_pythagorasCeq12(t *testing.T) {
	setup()
	defer teardown()
	testPythForSpecificC(t, 12, 20, 0)
}

func Test_pythagorasCeq25(t *testing.T) {
	setup()
	defer teardown()
	testPythForSpecificC(t, 25, 25, 2)
}

//func Test_pythagorasCeq52Multi(t *testing.T) {
//	setup()
//	defer teardown()
//	testPythForSpecificC(t, 52, 52, 1)
//}

//func Test_pythagorasCeq100Multi(t *testing.T) {
//	setup()
//	defer teardown()
//	testPythForSpecificC(t, 100, 100, 2)
//}
