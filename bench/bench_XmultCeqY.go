package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"testing"
)

func main() {
	bench_XmultCeqY()
}

// the driver for everything benching IntVar
func bench_XmultCeqY() {
	benchd(b_XmultCeqY1, bc{"name": "XmultCeqY", "size": "1"})
}

func b_XmultCeqY1(b *testing.B) { b_XmultCeqY(b) }

func b_XmultCeqY(b *testing.B) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		X := core.CreateIntVarFromTo("X", store, 100, 10000)
		Y := core.CreateIntVarFromTo("Y", store, 100, 1000)
		c := 10
		store.AddPropagator(explicit.CreateXmultCeqY(X, c, Y))
		store.IsConsistent()
	}
}
