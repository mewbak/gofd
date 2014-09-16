package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"testing"
)

func main() {
	bench_GC1XeqC2Y()
}

// the driver for everything benching IntVar
func bench_GC1XeqC2Y() {
	benchd(bGC1XeqC2Y1, bc{"name": "GC1XeqC2Y", "size": "1"})
}

func bGC1XeqC2Y1(b *testing.B) { bGC1XeqC2Y(b) }

func bGC1XeqC2Y(b *testing.B) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		X := core.CreateIntVarFromTo("X", store, 100, 10000)
		Y := core.CreateIntVarFromTo("Y", store, 100, 1000)
		c1, c2 := 1, 10
		store.AddPropagator(explicit.CreateC1XeqC2YBounds(c1, X, c2, Y))
		store.IsConsistent()
	}
}
